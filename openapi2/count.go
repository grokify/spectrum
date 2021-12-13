package openapi2

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/grokify/gocharts/data/histogram"
	"github.com/grokify/mogo/encoding/csvutil"
	"github.com/grokify/mogo/type/stringsutil"
)

func CountEndpointsByTag(spec Specification, tagsFilter []string) *histogram.HistogramSet {
	tagsFilter = stringsutil.SliceCondenseSpace(tagsFilter, true, true)
	hist := histogram.NewHistogramSet("endpoints by tag")
	for url, path := range spec.Paths {
		hist = countEndpointByTag(hist, tagsFilter, url, http.MethodGet, path.Get)
		hist = countEndpointByTag(hist, tagsFilter, url, http.MethodPatch, path.Patch)
		hist = countEndpointByTag(hist, tagsFilter, url, http.MethodPut, path.Put)
		hist = countEndpointByTag(hist, tagsFilter, url, http.MethodPost, path.Post)
		hist = countEndpointByTag(hist, tagsFilter, url, http.MethodDelete, path.Delete)
	}
	return hist
}

func countEndpointByTag(hist *histogram.HistogramSet, tagsFilter []string, url string, method string, ep *Endpoint) *histogram.HistogramSet {
	if ep == nil {
		return hist
	}
	method = strings.ToUpper(strings.TrimSpace(method))
	url = strings.TrimSpace(url)
	endpoint := method + " " + url
	for _, tag := range ep.Tags {
		tag = strings.TrimSpace(tag)
		add := true
		if len(tagsFilter) > 0 { // have tagsFilter
			add = false
			for _, try := range tagsFilter {
				if tag == try {
					add = true
				}
			}
		}
		if !add {
			continue
		}
		if len(tag) > 0 {
			hist.Add(tag, endpoint, 1)
		}
	}
	return hist
}

func WriteEndpointCountCSV(filename string, hset histogram.HistogramSet) error {
	writer, file, err := csvutil.NewWriterFile(filename)
	if err != nil {
		return err
	}
	//defer file.Close()
	//defer writer.Close()
	header := []string{"Tag", "Tag Endpoint Count", "Method", "Path"}
	err = writer.Write(header)
	if err != nil {
		return err
	}
	for tagName, hist := range hset.HistogramMap {
		hist.Inflate()
		for endpoint := range hist.Bins {
			parts := strings.Split(endpoint, " ")
			if len(parts) >= 2 {
				row := []string{
					tagName,
					strconv.Itoa(int(hist.BinCount)),
					strings.ToUpper(parts[0]),
					strings.Join(parts[1:], " ")}
				err := writer.Write(row)
				if err != nil {
					return err
				}
			}
		}
	}
	writer.Flush()
	err = writer.Error()
	if err != nil {
		return err
	}
	return file.Close()
}

// EndpointCount returns a count of the endpoints for a specification.
func EndpointCount(spec Specification) int {
	endpoints := map[string]int{}
	for url, path := range spec.Paths {
		url = strings.TrimSpace(url)
		if path.Get != nil && !path.Get.IsEmpty() {
			endpoints[fmt.Sprintf("%s %s", http.MethodGet, url)] = 1
		}
		if path.Patch != nil && !path.Patch.IsEmpty() {
			endpoints[fmt.Sprintf("%s %s", http.MethodPatch, url)] = 1
		}
		if path.Post != nil && !path.Post.IsEmpty() {
			endpoints[fmt.Sprintf("%s %s", http.MethodPost, url)] = 1
		}
		if path.Put != nil && !path.Put.IsEmpty() {
			endpoints[fmt.Sprintf("%s %s", http.MethodPut, url)] = 1
		}
		if path.Delete != nil && !path.Delete.IsEmpty() {
			endpoints[fmt.Sprintf("%s %s", http.MethodDelete, url)] = 1
		}
	}
	return len(endpoints)
}
