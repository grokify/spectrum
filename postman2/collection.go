package postman2

import (
	"encoding/json"
	"strings"
)

type Collection struct {
	Info  CollectionInfo `json:"info"`
	Item  []*Item        `json:"item"`
	Event []Event        `json:"event,omitempty"`
}

func NewCollectionFromBytes(data []byte) (Collection, error) {
	pman := Collection{}
	err := json.Unmarshal(data, &pman)
	return pman, err
}

func (col *Collection) GetOrNewFolder(folderName string) *Item {
	for _, folder := range col.Item {
		if folder.Name == folderName {
			return folder
		}
	}
	folder := &Item{Name: folderName, Item: []*Item{}}
	col.Item = append(col.Item, folder)
	return folder
}

func (col *Collection) SetFolder(newFolder *Item) {
	if newFolder == nil || len(strings.TrimSpace(newFolder.Name)) == 0 {
		return
	}
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
				folder.Item[j].Request.URL = &url
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

// Item can represent a folder or an API
type Item struct {
	Name        string   `json:"name,omitempty"`                 // Folder,Operation
	Description string   `json:"description,omitempty"`          // Folder
	Item        []*Item  `json:"item,omitempty"`                 // Folder
	IsSubFolder bool     `json:"_postman_isSubFolder,omitempty"` // Folder
	Event       []Event  `json:"event,omitempty"`                // Operation
	Request     *Request `json:"request,omitempty"`              // Operation
}

func (item *Item) UpsertSubItem(newItem *Item) {
	if newItem == nil || len(strings.TrimSpace(newItem.Name)) == 0 {
		return
	}
	for i, itemTry := range item.Item {
		if itemTry.Name == newItem.Name {
			item.Item[i] = newItem
			return
		}
	}
	item.Item = append(item.Item, newItem)
	return
}

type Event struct {
	Listen string `json:"listen"`
	Script Script `json:"script"`
}

type Script struct {
	Type string   `json:"type,omitempty"`
	Exec []string `json:"exec,omitmpety"`
}
