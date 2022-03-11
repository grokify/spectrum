package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/buaazp/fasthttprouter"
	"github.com/grokify/gohttp/anyhttp"
	"github.com/grokify/gohttp/httpsimple"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/net/httputilmore"
	"github.com/grokify/mogo/strconv/strconvutil"
	"github.com/grokify/mogo/type/stringsutil"
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
		Engine: stringsutil.FirstNonEmpty(os.Getenv("ENGINE"), "nethttp")}
}

func (svr *Server) HandleAPIRegistryNetHTTP(res http.ResponseWriter, req *http.Request) {
	log.Debug().Msg("FUNC_HandleNetHTTP__BEGIN")
	svr.HandleAPIRegistryAnyEngine(anyhttp.NewResReqNetHTTP(res, req))
}

func (svr *Server) HandleAPIRegistryFastHTTP(ctx *fasthttp.RequestCtx) {
	log.Debug().Msg("HANDLE_FastHTTP")
	svr.HandleAPIRegistryAnyEngine(anyhttp.NewResReqFastHTTP(ctx))
}

func (svr *Server) HandleAPIRegistryAnyEngine(aRes anyhttp.Response, aReq anyhttp.Request) {
	log.Debug().Msg("FUNC_HandleAnyEngine__BEGIN")
	aRes.SetContentType(httputilmore.ContentTypeTextHTMLUtf8)
	err := aReq.ParseForm()
	if err != nil {
		logutil.PrintErr(SetResponseError(aRes, err.Error()))
		return
	}

	specURL := strings.TrimSpace(aReq.QueryArgs().GetString("url"))
	log.Debug().
		Str("url", specURL).
		Msg("SpecURL")

	if len(specURL) == 0 {
		logutil.PrintErr(SetResponseError(aRes, "No OpenAPI 3.0 Spec URL"))
		return
	}
	resp, err := http.Get(specURL) // #nosec G107
	if err != nil || resp.StatusCode > 299 {
		logutil.PrintErr(SetResponseError(aRes, err.Error()))
		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logutil.PrintErr(SetResponseError(aRes, err.Error()))
		return
	}

	spec, err := openapi3.Parse(bytes)
	if err != nil {
		logutil.PrintErr(SetResponseError(aRes, err.Error()))
		return
	}

	oas3HTMLParams := openapi3html.PageParams{
		PageTitle:  spec.Info.Title,
		TableDomID: "specTable",
		Spec:       spec}

	oas3PageHTML := openapi3html.SpectrumUIPage(oas3HTMLParams)

	aRes.SetHeader(httputilmore.HeaderContentType, httputilmore.ContentTypeTextHTMLUtf8)
	_, err = aRes.SetBodyBytes([]byte(oas3PageHTML))
	logutil.PrintErr(err)
}

func SetResponseError(aRes anyhttp.Response, bodyText string) error {
	bodyText = `<!DOCTYPE html>
	<html>
	<body>
	<h1>` + bodyText + `</p></body></html>`

	_, err := aRes.SetBodyBytes([]byte(bodyText))
	return err
}

/*
const ErrorPage = `<!DOCTYPE html>
<html>
<h1>Error</h1>
<html>`
*/

func (svr Server) PortInt() int                       { return svr.Port }
func (svr Server) HttpEngine() string                 { return svr.Engine }
func (svr Server) RouterFast() *fasthttprouter.Router { return nil }

func (svr Server) Router() http.Handler {
	/*
		mux := mux.NewRouter()
		mux.HandleFunc("/ping", http.HandlerFunc(httpsimple.HandleTestNetHTTP))
		mux.HandleFunc("/ping/", http.HandlerFunc(httpsimple.HandleTestNetHTTP))
		mux.HandleFunc("/", http.HandlerFunc(svr.HandleAPIRegistryNetHTTP))
		return mux
	*/
	mux := http.NewServeMux()
	mux.HandleFunc("/", http.HandlerFunc(svr.HandleAPIRegistryNetHTTP))
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
