package postman2

import (
	"encoding/json"
	"regexp"
	"strings"
)

type Collection struct {
	Info CollectionInfo `json:"info"`
	Item []FolderItem   `json:"item"`
	Auth Auth           `json:"auth,omitempty"`
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
			if len(strings.TrimSpace(api.Request.URL.Raw)) > 0 {
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

type Auth struct {
	Type   string `json:"type,omitempty"`
	Bearer Bearer `json:"bearer,omitempty"`
}

type Bearer struct {
	Token string `json:"token,omitempty"`
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
	Query    []QueryParam      `json:"query,omitempty"`
}

type URLVariable struct {
	Value interface{} `json:"value,omitempty"`
	ID    string      `json:"id,omitempty"`
}

type QueryParam struct {
	Key         string `json:"key,omitempty"`
	Value       string `json:"value,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
	Description string `json:"description,omitempty"`
}

func NewURL(rawURL string) URL {
	rawURL = strings.TrimSpace(rawURL)
	url := URL{Raw: rawURL, Variable: []URLVariable{}}
	rx := regexp.MustCompile(`^([a-z]+)://([^/]+)/(.*)$`)
	rs := rx.FindAllStringSubmatch(rawURL, -1)

	if len(rs) > 0 {
		for _, m := range rs {
			url.Protocol = m[1]
			hostname := m[2]
			path := m[3]
			hostnameParts := strings.Split(hostname, ".")
			url.Host = hostnameParts

			pathParts := strings.Split(path, "/")
			url.Path = pathParts
		}
	}

	return url
}

func (url *URL) AddVariable(key string, value interface{}) {
	variable := URLVariable{ID: key, Value: value}
	url.Variable = append(url.Variable, variable)
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
