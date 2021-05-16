package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/buaazp/fasthttprouter"
	"github.com/grokify/simplego/net/anyhttp"
	"github.com/grokify/simplego/net/http/httpsimple"
	"github.com/grokify/simplego/net/httputilmore"
	"github.com/grokify/simplego/strconv/strconvutil"
	"github.com/grokify/simplego/type/stringsutil"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3/openapi3html"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

type Server struct {
	Port   int
	Engine string
}

func NewServer() Server {
	return Server{
		Port:   strconvutil.AtoiOrDefault(os.Getenv("PORT"), 8080),
		Engine: stringsutil.OrDefault(os.Getenv("ENGINE"), "nethttp")}
}

func (svc *Server) HandleAPIRegistryNetHTTP(res http.ResponseWriter, req *http.Request) {
	log.Debug().Msg("FUNC_HandleNetHTTP__BEGIN")
	svc.HandleAPIRegistryAnyEngine(anyhttp.NewResReqNetHttp(res, req))
}

func (svc *Server) HandleAPIRegistryFastHTTP(ctx *fasthttp.RequestCtx) {
	log.Debug().Msg("HANDLE_FastHTTP")
	svc.HandleAPIRegistryAnyEngine(anyhttp.NewResReqFastHttp(ctx))
}

func (svc *Server) HandleAPIRegistryAnyEngine(aRes anyhttp.Response, aReq anyhttp.Request) {
	log.Debug().Msg("FUNC_HandleAnyEngine__BEGIN")
	aRes.SetContentType(httputilmore.ContentTypeTextHtmlUtf8)
	err := aReq.ParseForm()
	if err != nil {
		SetResponseError(aRes, err.Error())
		return
	}

	specURL := strings.TrimSpace(aReq.QueryArgs().GetString("url"))
	log.Debug().
		Str("url", specURL).
		Msg("SpecURL")

	if len(specURL) == 0 {
		SetResponseError(aRes, "No OpenAPI 3.0 Spec URL")
		return
	}
	resp, err := http.Get(specURL)
	if err != nil || resp.StatusCode > 299 {
		SetResponseError(aRes, err.Error())
		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		SetResponseError(aRes, err.Error())
		return
	}

	spec, err := openapi3.Parse(bytes)
	if err != nil {
		SetResponseError(aRes, err.Error())
		return
	}

	oas3HtmlParams := openapi3html.PageParams{
		PageTitle:  spec.Info.Title,
		TableDomID: "specTable",
		Spec:       spec}

	oas3PageHtml := openapi3html.SwaggmanUIPage(oas3HtmlParams)

	aRes.SetHeader(httputilmore.HeaderContentType, httputilmore.ContentTypeTextHtmlUtf8)
	aRes.SetBodyBytes([]byte(oas3PageHtml))
}

func SetResponseError(aRes anyhttp.Response, bodyText string) {
	bodyText = `<!DOCTYPE html>
	<html>
	<body>
	<h1>` + bodyText + `</p></body></html>`

	aRes.SetBodyBytes([]byte(bodyText))
}

const ErrorPage = `<!DOCTYPE html>
<html>
<h1>Error</h1>
<html>`

func (svr Server) PortInt() int                       { return svr.Port }
func (svr Server) HttpEngine() string                 { return svr.Engine }
func (svr Server) RouterFast() *fasthttprouter.Router { return nil }

func (svc Server) Router() http.Handler {
	/*
		mux := mux.NewRouter()
		mux.HandleFunc("/ping", http.HandlerFunc(httpsimple.HandleTestNetHTTP))
		mux.HandleFunc("/ping/", http.HandlerFunc(httpsimple.HandleTestNetHTTP))
		mux.HandleFunc("/", http.HandlerFunc(svc.HandleAPIRegistryNetHTTP))
		return mux
	*/
	mux := http.NewServeMux()
	mux.HandleFunc("/", http.HandlerFunc(svc.HandleAPIRegistryNetHTTP))
	mux.HandleFunc("/ping", http.HandlerFunc(httpsimple.HandleTestNetHTTP))
	mux.HandleFunc("/ping/", http.HandlerFunc(httpsimple.HandleTestNetHTTP))
	return mux
}

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	svr := NewServer()

	done := make(chan bool)
	go httpsimple.Serve(svr)
	log.Info().Int("port", svr.Port).Msg("Server listening")
	<-done
}
