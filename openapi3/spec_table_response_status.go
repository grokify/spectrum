package openapi3

import (
	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gocharts/v2/data/histogram"
)

func (sm *SpecMore) StatusCodesHistogram() *histogram.HistogramSets {
	hsets := histogram.NewHistogramSets("Response Codes by Endpoint")
	if sm.Spec == nil {
		return hsets
	}
	VisitOperations(sm.Spec, func(path, method string, op *oas3.Operation) {
		/*
			if op == nil ||
				op.Responses == nil ||
				len(op.Responses) == 0 { // getkin v0.121.0 to v0.122.0
				return
			}
		*/
		if op == nil || op.Responses == nil {
			return
		}
		responsesMap := op.Responses.Map() // getkin v0.121.0 to v0.122.0
		if len(responsesMap) == 0 {
			return
		}
		for responseStatusCode := range responsesMap {
			hsets.Add(path, method, responseStatusCode, 1, true)
		}
	})
	return hsets
}

func (sm *SpecMore) WriteFileXLSXOperationStatusCodes(filename string) error {
	hsets := sm.StatusCodesHistogram()
	// return hsets.WriteXLSXPivot(filename, hsets.Name, "Method", "Path", "", "")
	return hsets.WriteXLSXPivot(filename, histogram.TablePivotOpts{
		TableName:           hsets.Name,
		ColNameHistogramSet: "Method",
		ColNameHistogram:    "Path",
	})
}
