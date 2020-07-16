package solr

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
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
	OptionMM                       = "mm"
	OptionBoost                    = "boost"
	OptionQueryFields              = "qf"
	OptionBoostQuery               = "bq"
	OptionBoostFunctions           = "bf"
	OptionUserFields               = "uf"
	ReturnTypeJSON                 = "json"
	QOperationOR                   = "OR"
	QOperationAND                  = "AND"
	DefTypeDisMax        DefType   = "dismax"
	DefTypeEDisMax       DefType   = "edismax"
	DefTypeStandard      DefType   = "lucene"
	DebugTypeQuery       DebugType = "query"
	DebugTypeTiming      DebugType = "timing"
	DebugTypeResults     DebugType = "results"
	DebugTypeAll         DebugType = "all"
)

// DebugType is used to restrict the available debug types for a
// `/search` request
type DebugType string

func (dt DebugType) String() string {
	return string(dt)
}

func (dt DebugType) isValid() bool {
	return !(dt != DebugTypeQuery && dt != DebugTypeTiming && dt != DebugTypeResults && dt != DebugTypeAll)
}

// DefType is used to restrict the available defTypes for a
// `/search` request
type DefType string

func (dt DefType) String() string {
	return string(dt)
}

func (dt DefType) isValid() bool {
	return !(dt != DefTypeDisMax && dt != DefTypeEDisMax && dt != DefTypeStandard)
}

// Returned validation errors
var (
	ErrInvalidDefType   = errors.New("invalid defType, please use one of the provided ones")
	ErrInvalidDebugType = errors.New("invalid debugType, please use one of the provided ones")
)

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

// ReadOptions contains options for read actions. Those include:
// Debug: Sets the type of debugging for the request
// DefType: Sets the type of query parse to use (default: lucene)
// Rows: Sets the number of rows to return
type ReadOptions struct {
	Debug   DebugType
	DefType DefType
	Rows    int
}

// Query represents the query parameters of a search. It provides
// helper methods for most of the available solr query params.
type Query struct {
	Q      []string
	params url.Values
}

// NewQuery returns an initialized Query. It accepts as options a result
// rows limit and a debug type. It sets by default the return type
// to JSON, as it is the only type supported by this library.
func NewQuery(opts *ReadOptions) *Query {
	nq := &Query{}
	nq.params = make(url.Values)
	if opts != nil {
		if opts.Debug != "" && opts.Debug.isValid() {
			nq.params.Set(OptionDebug, opts.Debug.String())
		}
		if opts.DefType != "" && opts.DefType.isValid() {
			nq.params.Set(OptionDefType, opts.DefType.String())
		}
		if opts.Rows > 0 {
			sv := strconv.Itoa(opts.Rows)
			nq.params.Set(OptionRows, sv)
		}
	}
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

// AddQuery adds a key-value pair to the Q parameter. Warning:
// Using this will overwrite any call to the `SetQuery`
// method. For complex logic use that instead
func (q *Query) AddQuery(key, value string) {
	q.Q = append(q.Q, fmt.Sprintf("%s:%s", key, value))
}

// SetQuery sets the Q parameter of the query.
func (q *Query) SetQuery(value string) {
	q.params.Set(OptionQ, value)
}

// SetOperationAND sets the operation for the Q parameter
// to AND (only when using `AddQuery`)
func (q *Query) SetOperationAND() {
	q.params.Set(OptionQOperation, QOperationAND)
}

// SetOperationOR sets the operation for the Q parameter
// to OR (only when using `AddQuery`)
func (q *Query) SetOperationOR() {
	q.params.Set(OptionQOperation, QOperationOR)
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

// AddField adds the given field to the returned field list.
// More info:
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

// SetSort sets the way the results are sorted. It should be formatted using the following
// protocol "<field name> <direction>, <field name> <direction>,...​"
// More info:
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
	if len(q.Q) > 0 {
		q.params.Set(OptionQ, strings.Join(q.Q, ", "))
	}
	q.params.Set(OptionWT, ReturnTypeJSON)
	return q.params.Encode()
}

// CollapseParams are the available params that can be set when using
// the Collapsing Query Parser
type CollapseParams struct {
	Field      string
	Min        string
	Max        string
	Sort       string
	NullPolicy *NullPolicy
	Hint       *Hint
	Size       string
}

// NullPolicy determines the policy when the collapsing field
// value is null on the document
type NullPolicy string

func (p NullPolicy) String() string {
	return string(p)
}

// Hint represents the Collapse hint param
type Hint string

func (h Hint) String() string {
	return string(h)
}

// Constants to secure proper NullPolicy & Hint usage
const (
	NullPolicyIgnore   NullPolicy = "ignore"
	NullPolicyExpand   NullPolicy = "expand"
	NullPolicyCollapse NullPolicy = "collapse"
	HintTopFC          Hint       = "top_fc"
)

// Possible errors returned from improper use of the Collapsing Query Parser
var (
	ErrParamsRequired    = errors.New("Param field is required for the CollapsingQParser")
	ErrTooManyParams     = errors.New("Only one of Max, Min or Sort may be populated")
	ErrInvalidNullPolicy = errors.New("invalid null policy, please use one of the provided")
	ErrInvalidHint       = errors.New("invalid hint, please use one of the provided")
)

func paramFormat(k, v string) string {
	return fmt.Sprintf("%s=%s", k, v)
}

func (p *CollapseParams) format() (string, error) {
	if p == nil {
		return "", ErrParamsRequired
	}
	if p.Field == "" {
		return "", ErrParamsRequired
	}

	params := []string{paramFormat("field", p.Field)}

	var c int
	if p.Max != "" {
		params = append(params, paramFormat("max", p.Max))
		c++
	}
	if p.Min != "" {
		params = append(params, paramFormat("min", p.Min))
		c++
	}
	if p.Sort != "" {
		params = append(params, paramFormat("sort", p.Sort))
		c++
	}
	if c > 1 {
		return "", ErrTooManyParams
	}

	if p.NullPolicy != nil {
		if *p.NullPolicy != NullPolicyIgnore && *p.NullPolicy != NullPolicyCollapse && *p.NullPolicy != NullPolicyExpand {
			return "", ErrInvalidNullPolicy
		}
		params = append(params, paramFormat("nullPolicy", p.NullPolicy.String()))
	}

	if p.Hint != nil {
		if *p.Hint != HintTopFC {
			return "", ErrInvalidHint
		}
		params = append(params, paramFormat("hint", p.Hint.String()))
	}

	if p.Size != "" {
		params = append(params, paramFormat("size", p.Size))
	}

	return fmt.Sprintf("{!collapse %s}", strings.Join(params, " ")), nil
}

// Collapse sets the Collapsing Query Parser post filter that groups
// documents according to the given parameters.
// More Info:
// https://lucene.apache.org/solr/guide/8_5/collapse-and-expand-results.html#collapsing-query-parser
func (q *Query) Collapse(params *CollapseParams) error {
	collapseString, err := params.format()
	if err != nil {
		return err
	}
	q.params.Add(OptionFilter, collapseString)
	return nil
}

// ExpandOptions are the available options to set for the expand component
type ExpandOptions struct {
	Sort string
	Rows int
	Q    string
	FQ   string
}

// Expand sets the parameter than returns an expand component used to expand the groups
// that were collapsed by the Collapsing Query Parser. The optional params override
// the original query values
// More info:
// https://lucene.apache.org/solr/guide/8_5/collapse-and-expand-results.html#expand-component
func (q *Query) Expand(opts *ExpandOptions) {
	q.params.Add("expand", "true")
	if opts != nil {
		if opts.Sort != "" {
			q.params.Add("expand.sort", opts.Sort)
		}
		if opts.Q != "" {
			q.params.Add("expand.q", opts.Q)
		}
		if opts.FQ != "" {
			q.params.Add("expand.fq", opts.FQ)
		}
		if opts.Rows > 0 {
			rv := strconv.Itoa(opts.Rows)
			q.params.Add("expand.rows", rv)
		}
	}
}

// SetQueryFields sets the fields to search (DisMax & eDisMax only)
// More info:
// https://lucene.apache.org/solr/guide/8_5/the-dismax-query-parser.html#qf-query-fields-parameter
func (q *Query) SetQueryFields(fields []string) {
	fieldsStr := strings.Join(fields, " ")
	q.params.Set(OptionQueryFields, fieldsStr)
}

// BoostField is a helper function to properly format field boosting
func BoostField(field string, boost float64) string {
	return fmt.Sprintf("%s^%f", field, boost)
}

// SetMinimumShouldMatch sets the minimum params to match (DisMax & eDisMax only)
// More info:
// https://lucene.apache.org/solr/guide/8_5/the-dismax-query-parser.html#mm-minimum-should-match-parameter
func (q *Query) SetMinimumShouldMatch(value string) {
	q.params.Set(OptionMM, value)
}

// SetBoostFunctions sets the boost functions param (DisMax & eDisMax only)
// More info:
// https://lucene.apache.org/solr/guide/8_5/the-dismax-query-parser.html#bf-boost-functions-parameter
func (q *Query) SetBoostFunctions(value string) {
	q.params.Set(OptionBoostFunctions, value)
}

// SetBoostQuery sets the boost query param (DisMax & eDisMax only)
// More info:
// https://lucene.apache.org/solr/guide/8_5/the-dismax-query-parser.html#bq-boost-query-parameter
func (q *Query) SetBoostQuery(value string) {
	q.params.Set(OptionBoostQuery, value)
}

// SetBoost sets the boost param (eDisMax only)
// More info:
// https://lucene.apache.org/solr/guide/8_5/the-extended-dismax-query-parser.html#extended-dismax-parameters
func (q *Query) SetBoost(value string) {
	q.params.Set(OptionBoost, value)
}

// SetUserFields sets the fields a user is allowed to query (eDisMax only)
// More info:
// https://lucene.apache.org/solr/guide/8_5/the-extended-dismax-query-parser.html#extended-dismax-parameters
func (q *Query) SetUserFields(fields []string) {
	fieldsStr := strings.Join(fields, " ")
	q.params.Set(OptionUserFields, fieldsStr)
}
