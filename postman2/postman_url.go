package postman2

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/grokify/gotilla/net/urlutil"
)

// URL is the Postman URL used in the Postman 2.0 Collection Spec.
type URL struct {
	Raw      string            `json:"raw,omitempty"`
	Protocol string            `json:"protocol,omitempty"`
	Auth     map[string]string `json:"auth"`
	Host     []string          `json:"host,omitempty"`
	Path     []string          `json:"path,omitempty"`
	Variable []URLVariable     `json:"variable,omitempty"`
}

func (url *URL) IsRawOnly() bool {
	url.Protocol = strings.TrimSpace(url.Protocol)
	if len(url.Protocol) > 0 ||
		len(url.Host) > 0 ||
		len(url.Path) > 0 {
		return false
	}
	return true
}

type URLVariable struct {
	Value interface{} `json:"value,omitempty"`
	ID    string      `json:"id,omitempty"`
}

func NewURLForGoUrl(goUrl url.URL) URL {
	pmURL := URL{Variable: []URLVariable{}}
	goUrl.Scheme = strings.TrimSpace(goUrl.Scheme)
	goUrl.Host = strings.TrimSpace(goUrl.Host)
	goUrl.Path = strings.TrimSpace(goUrl.Path)
	urlParts := []string{}
	if len(goUrl.Host) > 0 {
		pmURL.Host = strings.Split(goUrl.Host, ".")
		urlParts = append(urlParts, goUrl.Host)
	}
	if len(goUrl.Path) > 0 {
		pmURL.Path = strings.Split(goUrl.Path, "/")
		urlParts = append(urlParts, goUrl.Path)
	}
	rawURL := strings.Join(urlParts, "/")
	if len(goUrl.Scheme) > 0 {
		pmURL.Protocol = goUrl.Scheme
		rawURL = goUrl.Scheme + "://" + rawURL
	}
	pmURL.Raw = rawURL
	return pmURL
}

var simpleURLRx *regexp.Regexp = regexp.MustCompile(`^([a-z][0-9a-z]+)://([^/]+)/(.*)$`)

func NewURL(rawURL string) URL {
	rawURL = strings.TrimSpace(rawURL)
	pmURL := URL{Raw: rawURL, Variable: []URLVariable{}}
	rs1 := simpleURLRx.FindAllStringSubmatch(rawURL, -1)

	if len(rs1) > 0 {
		for _, m := range rs1 {
			pmURL.Protocol = m[1]
			hostname := m[2]
			path := m[3]
			pmURL.Host = strings.Split(hostname, ".")
			pmURL.Path = urlutil.SplitPath(path, true, true)
		}
	} else if strings.Index(rawURL, "{") == 0 {
		parts := urlutil.SplitPath(rawURL, true, true)
		if len(parts) > 0 {
			pmURL.Host = []string{parts[0]}
		}
		if len(parts) > 1 {
			pmURL.Path = parts[1:]
		}
	}

	return pmURL
}

// AddVariable adds a Postman Variable to the struct.
func (pmURL *URL) AddVariable(key string, value interface{}) {
	variable := URLVariable{ID: key, Value: value}
	pmURL.Variable = append(pmURL.Variable, variable)
}

const (
	apiUrlOasToPostmanVarMatch   string = `(^|[^\{])\{([^\/\{\}]+)\}([^\}]|$)`
	apiUrlOasToPostmanVarReplace string = "$1:$2$3"
)

var apiUrlOasToPostmanVarMatchRx *regexp.Regexp = regexp.MustCompile(
	apiUrlOasToPostmanVarMatch)

func ApiUrlOasToPostman(urlWithOasVars string) string {
	return apiUrlOasToPostmanVarMatchRx.ReplaceAllString(
		urlWithOasVars, apiUrlOasToPostmanVarReplace)
}
