package solr

import (
	"encoding/json"
	"time"
)

// Response represents the response from the solr server. It usually contains
// Header information, the response data or an error in case of erroneous
// response. Also it can contain Debug information when requested, a
// single document (in the case of realtimeGet) or just a status
// (in the case of the Ping request)
type Response struct {
	Header      *ResponseHeader          `json:"responseHeader"`
	Data        *ResponseData            `json:"response"`
	Error       *ResponseError           `json:"error"`
	Debug       *map[string]interface{}  `json:"debug"`
	Doc         *Doc                     `json:"doc"`
	Status      *string                  `json:"status"`
	Expanded    map[string]*ResponseData `json:"expanded"`
	FacetCounts *FacetCounts             `json:"facet_counts"`
	Grouped     map[string]*GroupField   `json:"grouped"`
}

// ResponseHeader is populated on every response from the solr server
// unless explicitly omitted. It contains the request status code
// the time it took as well as the params for the search query
// when applicable
type ResponseHeader struct {
	Status int64                   `json:"status"`
	QTime  int64                   `json:"QTime"`
	Params *map[string]interface{} `json:"params"`
}

// ResponseData is populated on a successful response from the solr
// server. It contains the number of documents found, the starting
// index (in case of a search) as well as the documents found
type ResponseData struct {
	NumFound int64     `json:"numFound"`
	Start    int64     `json:"start"`
	Docs     Docs      `json:"docs"`
	MaxScore *MaxScore `json:"maxScore"` //TODO(KK): Fix this
}

// MaxScore is used as a struct due to the fact that solr
// may return it as a float or as a string indicating
// "NaN"
type MaxScore struct {
	Valid bool
	Score float64
}

// UnmarshalJSON implements the unmarshaler interface
func (m *MaxScore) UnmarshalJSON(b []byte) error {
	var i interface{}
	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}
	switch v := i.(type) {
	case float64:
		m.Score = v
		m.Valid = true
	default:
		m.Valid = false
	}
	return nil
}

// ResponseError is populated in the event the response from the solr
// server is erroneous. It contains the status code, a message
// and some metadata about the error's class
type ResponseError struct {
	Code    int64    `json:"code"`
	Message string   `json:"msg"`
	Meta    []string `json:"metadata"`
}

func (r *ResponseError) Error() string {
	return r.Message
}

// Docs represents an array of doc
type Docs []*Doc

// Doc is essentialy a map[string]interface{}
type Doc map[string]interface{}

// ToBytes returs a byte slice to simplify unmarshaling to JSON
func (d Docs) ToBytes() ([]byte, error) {
	return interfaceToBytes(d)
}

// ToBytes returs a byte slice to simplify unmarshaling to JSON
func (d *Doc) ToBytes() ([]byte, error) {
	return interfaceToBytes(d)
}

// FacetCounts is populated whenever the query to solr includes facets.
// Each of the following attributes get populated depending on the
// actual facet query. 'Fields' attribute includes a helper to
// retrieve the facets in a string: float format.
type FacetCounts struct {
	Queries   map[string]int         `json:"facet_queries"`
	Fields    *FacetFields           `json:"facet_fields"`
	Dates     map[string]interface{} `json:"facet_dates"`
	Ranges    map[string]*Range      `json:"facet_ranges"`
	Intervals map[string]interface{} `json:"facet_intervals"`
	Heatmaps  map[string]interface{} `json:"facet_heatmaps"`
	Pivot     map[string][]*Pivot    `json:"facet_pivot"`
}

// FacetFields is the facet_field parameter which in Solr contains an array
// that alternates between string and numbers. In order to make this
// more Go-friendly it's using a custom unmarshaler and a getter
// that helps format the results in a map[string]float64.
type FacetFields struct {
	m map[string]map[string]float64
}

// Get returns the facets on the given field in a Go-friendly way.
func (f *FacetFields) Get(s string) map[string]float64 {
	return f.m[s]
}

// UnmarshalJSON implements the unmarshaler interface.
func (f *FacetFields) UnmarshalJSON(b []byte) error {
	f.m = make(map[string]map[string]float64)
	var temp map[string][]interface{}
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}

	for k, v := range temp {
		values := map[string]float64{}
		for i := 0; i < len(v); i += 2 {
			s, ok := v[i].(string)
			n, ok2 := v[i+1].(float64)
			if ok && ok2 {
				values[s] = n
			}
		}
		f.m[k] = values
	}

	return nil
}

// Pivot contains pivot faceting results.
// More info:
// https://lucene.apache.org/solr/guide/8_5/faceting.html#pivot-decision-tree-faceting
type Pivot struct {
	Field   string            `json:"field"`
	Value   interface{}       `json:"value"`
	Count   int               `json:"count"`
	Pivot   []*Pivot          `json:"pivot"`
	Stats   *Stats            `json:"stats"`
	Queries map[string]int    `json:"queries"`
	Ranges  map[string]*Range `json:"ranges"`
}

// Range contains range faceting results.
// More info:
// https://lucene.apache.org/solr/guide/8_5/faceting.html#range-faceting
type Range struct {
	Counts *FacetFields `json:"counts"`
	Gap    string       `json:"gap"`
	Start  time.Time    `json:"start"`
	End    time.Time    `json:"end"`
}

// Stats containts the results of stats when requested during pivot faceting.
type Stats struct {
	Fields map[string]interface{} `json:"stats_fields"`
}

// GroupField is populated whenever the query to solr includes grouping.
// The response contains the total matches (of docs), the number of
// groups (if requested) and the groups.
type GroupField struct {
	Matches        int      `json:"matches"`
	NumberOfGroups int      `json:"ngroups"`
	Groups         []*Group `json:"groups"`
}

// Group contains a value and the list of documents that belong
// to the specific group.
type Group struct {
	Value   interface{}   `json:"groupValue"`
	DocList *ResponseData `json:"doclist"`
}
