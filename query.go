package solr

import (
	"fmt"
	"net/url"
	"strconv"
)

// Query Options and other constants
const (
	QueryOptionDebug                  = "debug"
	QueryOptionDefType                = "defType"
	QueryOptionQ                      = "q"
	QueryOptionQOperation             = "q.op"
	QueryOptionFilter                 = "fq"
	QueryOptionFieldList              = "fl"
	QueryOptionRows                   = "rows"
	QueryOptionStart                  = "start"
	QueryOptionSort                   = "sort"
	QueryOptionWT                     = "wt"
	QueryOptionCommit                 = "commit"
	QueryOptionOverwrite              = "overwrite"
	QueryOptionCommitWithin           = "commitWithin"
	ReturnTypeJSON                    = "json"
	DefTypeDisMax                     = "dismax"
	DefTypeEDisMax                    = "edismax"
	DefTypeStandard                   = "lucene"
	DebugTypeQuery          DebugType = "query"
)

type Query struct {
	Debug              string   `json:"debug"`
	FQ                 []string `json:"fq"`
	Sort               string   `json:"sort"`
	Start              string   `json:"start"`
	Rows               string   `json:"rows"`
	FieldList          []string `json:"fl"`
	DefaultSearchField string   `json:"df"`
	Raw                map[string]interface{}
	params             url.Values
}

type DebugType string

type QueryOptions struct {
	Debug *DebugType
	Rows  int
}

type WriteOptions struct {
	Commit       bool
	CommitWithin int64
	Overwrite    *bool
}

func (opts *WriteOptions) formatQueryFromOpts() url.Values {
	q := make(url.Values)
	if opts.Commit {
		q.Set(QueryOptionCommit, "true")
	}
	if opts.CommitWithin > 0 {
		q.Set(QueryOptionCommitWithin, strconv.FormatInt(opts.CommitWithin, 10))
	}
	if opts.Overwrite != nil {
		q.Set(QueryOptionOverwrite, strconv.FormatBool(*opts.Overwrite))
	}
	return q
}

// NewQuery returns a new Solr query
func NewQuery(opts *QueryOptions) *Query {
	nq := &Query{}
	nq.params = make(url.Values)
	if opts.Debug != nil {
		nq.params.Set(QueryOptionDebug, string(DebugTypeQuery))
	}
	if opts.Rows > 0 {
		sv := strconv.Itoa(opts.Rows)
		nq.params.Set(QueryOptionRows, sv)
	}
	nq.params.Set(QueryOptionWT, ReturnTypeJSON)
	return nq
}

// AddParam allows the addition of custom query parameters
func (q *Query) AddParam(key, value string) {
	q.params.Add(key, value)
}

// SetParam allows the setting of custom query parameters
func (q *Query) SetParam(key, value string) {
	q.params.Set(key, value)
}

// DelParam allows the deletion of query parameters
func (q *Query) DelParam(key string) {
	q.params.Del(key)
}

func (q *Query) SetQuery(value string) {
	q.params.Set(QueryOptionQ, value)
}

func (q *Query) SetOperation(value string) {
	q.params.Set(QueryOptionQOperation, value)
}

// AddFilter adds a key-value pair on which to filter the query
// https://lucene.apache.org/solr/guide/8_5/common-query-parameters.html#fq-filter-query-parameter
func (q *Query) AddFilter(key, value string) {
	q.params.Add(QueryOptionFilter, fmt.Sprintf("%s:%s", key, value))
}

// SetFilter gives the option to set a filter allowing for more complex logic instead
// of a basic key-value check
func (q *Query) SetFilter(value string) {
	q.params.Set(QueryOptionFilter, value)
}

// https://lucene.apache.org/solr/guide/8_5/common-query-parameters.html#fl-field-list-parameter
func (q *Query) AddField(value string) {
	q.params.Add(QueryOptionFieldList, value)
}

// https://lucene.apache.org/solr/guide/8_5/common-query-parameters.html#start-parameter
func (q *Query) SetStart(value int) {
	sv := strconv.Itoa(value)
	q.params.Set(QueryOptionStart, sv)
}

// <field name>+<direction>,<field name>+<direction>],…​
// https://lucene.apache.org/solr/guide/8_5/common-query-parameters.html#sort-parameter
func (q *Query) SetSort(value string) {
	q.params.Set(QueryOptionSort, value)
}

// // SetRows sets the amount of rows to be returned from the query overwritting the
// // default value lucene.apache.org/solr/guide/8_5/common-query-parameters.html#rows-parameter
// func (q *Query) SetRows(value int) {
// 	sv := strconv.Itoa(value)
// 	q.params.Set(QueryOptionRows, sv)
// }

func (q *Query) String() string {
	return q.params.Encode()
}
