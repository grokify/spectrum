package postman2

import (
	"encoding/json"
	"net/url"
	"regexp"
	"strings"
)

type Collection struct {
	Info CollectionInfo `json:"info"`
	Item []FolderItem   `json:"item"`
}

func NewCollectionFromBytes(data []byte) (Collection, error) {
	pman := Collection{}
	err := json.Unmarshal(data, &pman)
	return pman, err
}

func (col *Collection) GetOrNewFolder(folderName string) FolderItem {
	for _, folder := range col.Item {
		if folder.Name == folderName {
			return folder
		}
	}
	folder := FolderItem{
		Name: folderName}
	col.Item = append(col.Item, folder)
	return folder
}

func (col *Collection) SetFolder(newFolder FolderItem) {
	for i, folder := range col.Item {
		if newFolder.Name == folder.Name {
			col.Item[i] = newFolder
			return
		}
	}
	col.Item = append(col.Item, newFolder)
}

func (col *Collection) InflateRawURLs() {
	for _, folder := range col.Item {
		for j, api := range folder.Item {
			if api.Request.URL.IsRawOnly() &&
				len(strings.TrimSpace(api.Request.URL.Raw)) > 0 {
				url := NewURL(strings.TrimSpace(api.Request.URL.Raw))
				url.Auth = api.Request.URL.Auth
				url.Variable = api.Request.URL.Variable
				folder.Item[j].Request.URL = url
			}
		}
	}
}

type CollectionInfo struct {
	Name        string `json:"name,omitempty"`
	PostmanID   string `json:"_postman_id,omitempty"`
	Description string `json:"description,omitempty"`
	Schema      string `json:"schema,omitempty"`
}

type FolderItem struct {
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Item        []APIItem `json:"item,omitempty"`
}

type APIItem struct {
	Name    string  `json:"name,omitempty"`
	Event   []Event `json:"event,omitempty"`
	Request Request `json:"request,omitempty"`
}

type Event struct {
	Listen string `json:"listen"`
	Script Script `json:"script"`
}

type Script struct {
	Type string   `json:"type,omitempty"`
	Exec []string `json:"exec,omitmpety"`
}

type Request struct {
	URL         URL         `json:"url,omitempty"`
	Method      string      `json:"method,omitempty"`
	Header      []Header    `json:"header,omitempty"`
	Body        RequestBody `json:"body,omitempty"`
	Description string      `json:"description,omitempty"`
}

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

func NewURL(rawURL string) URL {
	rawURL = strings.TrimSpace(rawURL)
	pmURL := URL{Raw: rawURL, Variable: []URLVariable{}}
	rx := regexp.MustCompile(`^([a-z][0-9a-z]+)://([^/]+)/(.*)$`)
	rs := rx.FindAllStringSubmatch(rawURL, -1)

	if len(rs) > 0 {
		for _, m := range rs {
			pmURL.Protocol = m[1]
			hostname := m[2]
			path := m[3]
			pmURL.Host = strings.Split(hostname, ".")
			pmURL.Path = strings.Split(path, "/")
		}
	}

	return pmURL
}

func (pmURL *URL) AddVariable(key string, value interface{}) {
	variable := URLVariable{ID: key, Value: value}
	pmURL.Variable = append(pmURL.Variable, variable)
}

type Header struct {
	Key         string `json:"key,omitempty"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
}

type RequestBody struct {
	Mode       string            `json:"mode,omitempty"` // `raw`, `urlencoded`, `formdata`,`file`,`graphql`
	Raw        string            `json:"raw,omitempty"`
	URLEncoded []URLEncodedParam `json:"urlencoded,omitempty"`
}

type URLEncodedParam struct {
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
	Type    string `json:"type,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}
