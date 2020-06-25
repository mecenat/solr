package solr

import (
	"fmt"
	"net/url"
	"strconv"
)

// Query Options and other constants
const (
	OptionDebug                    = "debug"
	OptionDefType                  = "defType"
	OptionQ                        = "q"
	OptionQOperation               = "q.op"
	OptionFilter                   = "fq"
	OptionFieldList                = "fl"
	OptionRows                     = "rows"
	OptionStart                    = "start"
	OptionSort                     = "sort"
	OptionWT                       = "wt"
	OptionCommit                   = "commit"
	OptionOverwrite                = "overwrite"
	OptionCommitWithin             = "commitWithin"
	OptionWaitSearcher             = "waitSearcher"
	OptionMaxSegments              = "maxSegments"
	OptionExpungeDeletes           = "expungeDeletes"
	ReturnTypeJSON                 = "json"
	DefTypeDisMax        DefType   = "dismax"
	DefTypeEDisMax       DefType   = "edismax"
	DefTypeStandard      DefType   = "lucene"
	DebugTypeQuery       DebugType = "query"
)

// Query represents the query parameters of a search. It provides
// helper methods for most of the available solr query params.
type Query struct {
	Debug string `json:"debug"`
	// FQ                 []string `json:"fq"`
	Sort               string   `json:"sort"`
	Start              string   `json:"start"`
	Rows               string   `json:"rows"`
	FieldList          []string `json:"fl"`
	DefaultSearchField string   `json:"df"`
	Raw                map[string]interface{}
	params             url.Values
}

type DebugType string

type DefType string

type ReadOptions struct {
	Debug *DebugType
	Rows  int
}

// WriteOptions contains options for write actions. Those include:
// Commit: Autocommit all changes alongside the current request
// CommitWithin: Autocommit all changes after the specified
//     time (in miliseconds)
// AllowDuplicate: Allows uniqueKey duplication
type WriteOptions struct {
	Commit         bool
	CommitWithin   int64
	AllowDuplicate bool
}

func (opts *WriteOptions) formatQueryFromOpts() url.Values {
	if opts == nil {
		return nil
	}

	q := make(url.Values)
	if opts.Commit {
		q.Set(OptionCommit, "true")
	}
	if opts.CommitWithin > 0 {
		q.Set(OptionCommitWithin, strconv.FormatInt(opts.CommitWithin, 10))
	}
	if opts.AllowDuplicate {
		q.Set(OptionOverwrite, "false")
	}
	return q
}

// NewQuery returns an initialized Query. It accepts as options a result
// rows limit and a debug type. It sets by default the return type
// to JSON, as it is the only type supported by this library.
func NewQuery(opts *ReadOptions) *Query {
	nq := &Query{}
	nq.params = make(url.Values)
	if opts.Debug != nil {
		nq.params.Set(OptionDebug, string(DebugTypeQuery))
	}
	if opts.Rows > 0 {
		sv := strconv.Itoa(opts.Rows)
		nq.params.Set(OptionRows, sv)
	}
	nq.params.Set(OptionWT, ReturnTypeJSON)
	return nq
}

// AddParam allows the addition of custom query parameters.
func (q *Query) AddParam(key, value string) {
	q.params.Add(key, value)
}

// SetParam allows the setting of custom query parameters.
func (q *Query) SetParam(key, value string) {
	q.params.Set(key, value)
}

// DelParam allows the deletion of query parameters.
func (q *Query) DelParam(key string) {
	q.params.Del(key)
}

func (q *Query) SetQuery(value string) {
	q.params.Set(OptionQ, value)
}

func (q *Query) SetOperation(value string) {
	q.params.Set(OptionQOperation, value)
}

// AddFilter adds a key-value pair on which to filter the query.
// More info:
// https://lucene.apache.org/solr/guide/8_5/common-query-parameters.html#fq-filter-query-parameter
func (q *Query) AddFilter(key, value string) {
	q.params.Add(OptionFilter, fmt.Sprintf("%s:%s", key, value))
}

// SetFilter gives the option to set a filter allowing for more complex logic instead
// of a basic key-value check.
func (q *Query) SetFilter(value string) {
	q.params.Set(OptionFilter, value)
}

// https://lucene.apache.org/solr/guide/8_5/common-query-parameters.html#fl-field-list-parameter
func (q *Query) AddField(value string) {
	q.params.Add(OptionFieldList, value)
}

// SetStart enables setting the starting index for a search query. It can be used when
// the available results are more than the rows returned to fetch the remainder rows.
// More info:
// https://lucene.apache.org/solr/guide/8_5/common-query-parameters.html#start-parameter
func (q *Query) SetStart(value int) {
	sv := strconv.Itoa(value)
	q.params.Set(OptionStart, sv)
}

// <field name>+<direction>,<field name>+<direction>],…​
// https://lucene.apache.org/solr/guide/8_5/common-query-parameters.html#sort-parameter
func (q *Query) SetSort(value string) {
	q.params.Set(OptionSort, value)
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
