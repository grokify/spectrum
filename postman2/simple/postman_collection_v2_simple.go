package simple

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/grokify/simplego/net/httputilmore"
	"github.com/grokify/spectrum/postman2"
	"github.com/pkg/errors"
)

type Collection struct {
	Info postman2.CollectionInfo `json:"info"`
	Item []*Item                 `json:"item"`
}

func NewCollectionFromBytes(data []byte) (Collection, error) {
	pman := Collection{}
	err := json.Unmarshal(data, &pman)
	if err != nil {
		err = errors.Wrap(err, "spectrum.postman2.simple.NewCollectionFromBytes << json.Unmarshal")
	}
	return pman, err
}

func NewCanonicalCollectionFromBytes(data []byte) (postman2.Collection, error) {
	collection, errTry := postman2.NewCollectionFromBytes(data)
	if errTry == nil {
		collection.Inflate()
		return collection, nil
	}
	simpleCollection, err := NewCollectionFromBytes(data)
	if err != nil {
		err = errors.Wrap(errTry, err.Error())
		err = errors.Wrap(err, "spectrum.postman2.simple.NewCanonicalCollectionFromBytes << NewCollectionFromBytes")
		return collection, err
	}
	collection = simpleCollection.ToCanonical()
	collection.Inflate()
	return collection, nil
}

//noinspection ALL
func ReadCanonicalCollection(filepath string) (postman2.Collection, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		err = errors.Wrap(err, "spectrum.postman2.ReadCanonicalCollection << ioutil.ReadFile")
		return postman2.Collection{}, err
	}
	return NewCanonicalCollectionFromBytes(bytes)
}

func (col *Collection) ToCanonical() postman2.Collection {
	cCollection := postman2.Collection{
		Info: col.Info,
		Item: []*postman2.Item{}}
	for _, folder := range col.Item {
		cCollection.Item = append(cCollection.Item, folder.ToCanonical())
	}
	return cCollection
}

type Item struct {
	Name        string           `json:"name,omitempty"`        // Folder,API
	Description string           `json:"description,omitempty"` // Folder
	Item        []*Item          `json:"item,omitempty"`        // Folder
	Event       []postman2.Event `json:"event,omitempty"`       // API
	Request     Request          `json:"request,omitempty"`     // API
}

func (thisItem *Item) ToCanonical() *postman2.Item {
	canRequest := thisItem.Request.ToCanonical()
	canItem := &postman2.Item{
		Name:    thisItem.Name,
		Item:    []*postman2.Item{},
		Event:   thisItem.Event,
		Request: &canRequest}
	thisItem.Description = strings.TrimSpace(thisItem.Description)
	if len(thisItem.Description) > 0 {
		canItem.Description = &postman2.Description{
			Content: thisItem.Description,
			Type:    httputilmore.ContentTypeTextMarkdown}
	}
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
	canReq := apiItem.Request.ToCanonical()
	return postman2.Item{
		Name:    apiItem.Name,
		Event:   apiItem.Event,
		Request: &canReq}
}

type Request struct {
	URL         string               `json:"url,omitempty"`
	Method      string               `json:"method,omitempty"`
	Header      []postman2.Header    `json:"header,omitempty"`
	Body        postman2.RequestBody `json:"body,omitempty"`
	Description string               `json:"description,omitempty"`
}

func (req *Request) ToCanonical() postman2.Request {
	pmUrl := postman2.NewURL(req.URL)
	return postman2.Request{
		URL:         &pmUrl,
		Method:      req.Method,
		Header:      req.Header,
		Body:        &req.Body,
		Description: req.Description}
}
