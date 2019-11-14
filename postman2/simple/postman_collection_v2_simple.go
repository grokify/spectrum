package simple

import (
	"encoding/json"
	"io/ioutil"

	"github.com/grokify/swaggman/postman2"
)

type Collection struct {
	Info postman2.CollectionInfo `json:"info"`
	Item []Item                  `json:"item"`
}

func NewCollectionFromBytes(data []byte) (Collection, error) {
	pman := Collection{}
	err := json.Unmarshal(data, &pman)
	return pman, err
}

func NewCanonicalCollectionFromBytes(data []byte) (postman2.Collection, error) {
	cPman, err := postman2.NewCollectionFromBytes(data)
	if err == nil {
		cPman.InflateRawURLs()
		return cPman, nil
	}
	sPman, err := NewCollectionFromBytes(data)
	if err != nil {
		return cPman, err
	}
	cPman = sPman.ToCanonical()
	return cPman, nil
}

func ReadCanonicalCollection(filepath string) (postman2.Collection, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return postman2.Collection{}, err
	}
	return NewCanonicalCollectionFromBytes(bytes)
}

func (col *Collection) ToCanonical() postman2.Collection {
	cCollection := postman2.Collection{
		Info: col.Info,
		Item: []postman2.Item{}}
	for _, folder := range col.Item {
		cCollection.Item = append(cCollection.Item, folder.ToCanonical())
	}
	return cCollection
}

type Item struct {
	Name        string           `json:"name,omitempty"`        // Folder,API
	Description string           `json:"description,omitempty"` // Folder
	Item        []Item           `json:"item,omitempty"`        // Folder
	Event       []postman2.Event `json:"event,omitempty"`       // API
	Request     Request          `json:"request,omitempty"`     // API
}

func (thisItem *Item) ToCanonical() postman2.Item {
	canItem := postman2.Item{
		Name:        thisItem.Name,
		Description: thisItem.Description,
		Item:        []postman2.Item{},
		Event:       thisItem.Event,
		Request:     thisItem.Request.ToCanonical()}
	for _, subItem := range thisItem.Item {
		canItem.Item = append(canItem.Item, subItem.ToCanonical())
	}
	return canItem
}

type APIItem struct {
	Name    string           `json:"name,omitempty"`
	Event   []postman2.Event `json:"event,omitempty"`
	Request Request          `json:"request,omitempty"`
}

func (apiItem *APIItem) ToCanonical() postman2.Item {
	return postman2.Item{
		Name:    apiItem.Name,
		Event:   apiItem.Event,
		Request: apiItem.Request.ToCanonical()}
}

type Request struct {
	URL         string               `json:"url,omitempty"`
	Method      string               `json:"method,omitempty"`
	Header      []postman2.Header    `json:"header,omitempty"`
	Body        postman2.RequestBody `json:"body,omitempty"`
	Description string               `json:"description,omitempty"`
}

func (req *Request) ToCanonical() postman2.Request {
	return postman2.Request{
		URL:         postman2.NewURL(req.URL),
		Method:      req.Method,
		Header:      req.Header,
		Body:        req.Body,
		Description: req.Description}
}
