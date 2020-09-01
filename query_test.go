package solr

import (
	"testing"
)

func TestNewQuery(t *testing.T) {
	opts := &ReadOptions{
		Debug:   DebugTypeQuery,
		DefType: "test",
		Rows:    10,
	}
	q := NewQuery(opts)

	if q.params.Get("defType") == "test" {
		t.Fatal("did register invalid defType")
	}

	if q.params.Get("debug") != DebugTypeQuery.String() {
		t.Fatal("did not register valid debug type")
	}

	if q.params.Get("rows") != "10" {
		t.Fatal("did not register valid rows number")
	}

	opts2 := &ReadOptions{
		Debug:   "debug",
		DefType: DefTypeDisMax,
	}
	q = NewQuery(opts2)

	if q.params.Get("defType") != DefTypeDisMax.String() {
		t.Fatal("did not register valid defType")
	}

	if q.params.Get("debug") == "debug" {
		t.Fatal("did register invalid debug type")
	}
}

func TestAddQuery(t *testing.T) {
	q := NewQuery(nil)

	q.AddQuery("field", "value")
	q.AddQuery("", "value string")

	if len(q.q) != 2 {
		t.Fatalf("query length should be 2 is instead %d", len(q.q))
	}

	actual := q.String()
	expected := "q=field%3Avalue+OR+value+string&wt=json"
	if actual != expected {
		t.Fatalf("expected %s but got %s", expected, actual)
	}
}

func TestSetQuery(t *testing.T) {
	q := NewQuery(nil)
	q.SetQuery("key:value")
	if q.params.Get("q") == "" {
		t.Fatal("q param not registered")
	}
}

func TestDelQuery(t *testing.T) {
	q := NewQuery(nil)
	q.SetQuery("key:value")
	q.DelQuery()
	if q.params.Get("q") != "" {
		t.Fatal("q param not deleted properly")
	}
}

func TestAddFilter(t *testing.T) {
	q := NewQuery(nil)
	q.AddFilter("key", "value")
	if q.params.Get("fq") == "" {
		t.Fatal("fq param not registered")
	}
}

func TestSetFilter(t *testing.T) {
	q := NewQuery(nil)
	q.SetFilter("key:value")
	if q.params.Get("fq") == "" {
		t.Fatal("fq param not registered")
	}
}

func TestAddField(t *testing.T) {
	q := NewQuery(nil)
	q.AddField("key")
	if q.params.Get("fl") == "" {
		t.Fatal("fl param not registered")
	}
}

func TestSetStart(t *testing.T) {
	q := NewQuery(nil)
	q.SetStart(0)
	if q.params.Get("start") == "" {
		t.Fatal("start param not registered")
	}
}

func TestSetSort(t *testing.T) {
	q := NewQuery(nil)
	q.SetSort("key asc")
	if q.params.Get("sort") == "" {
		t.Fatal("sort param not registered")
	}
}

func TestCollapseNoParams(t *testing.T) {
	q := NewQuery(nil)
	err := q.Collapse(nil)
	if err == nil {
		t.Fatal("expected error but got nothing")
	}
}

func TestCollapseNoField(t *testing.T) {
	q := NewQuery(nil)
	params := &CollapseParams{
		Field: "",
	}
	err := q.Collapse(params)
	if err == nil {
		t.Fatal("expected error but got nothing")
	}
}
func TestCollapseBothMinMax(t *testing.T) {
	q := NewQuery(nil)
	params := &CollapseParams{
		Field: "field",
		Min:   "test",
		Max:   "test",
	}
	err := q.Collapse(params)
	if err == nil {
		t.Fatal("expected error but got nothing")
	}
}
func TestCollapseBothSortMin(t *testing.T) {
	q := NewQuery(nil)
	params := &CollapseParams{
		Field: "field",
		Min:   "test",
		Sort:  "test",
	}
	err := q.Collapse(params)
	if err == nil {
		t.Fatal("expected error but got nothing")
	}
}
func TestCollapseBothSortMax(t *testing.T) {
	q := NewQuery(nil)
	params := &CollapseParams{
		Field: "field",
		Max:   "test",
		Sort:  "test",
	}
	err := q.Collapse(params)
	if err == nil {
		t.Fatal("expected error but got nothing")
	}
}
func TestCollapseBadHint(t *testing.T) {
	q := NewQuery(nil)
	badHint := Hint("hint")
	params := &CollapseParams{
		Field: "field",
		Hint:  &badHint,
	}
	err := q.Collapse(params)
	if err == nil {
		t.Fatal("expected error but got nothing")
	}
}
func TestCollapseBadNullPolicy(t *testing.T) {
	q := NewQuery(nil)
	badNullPolicy := NullPolicy("nuls")
	params := &CollapseParams{
		Field:      "field",
		NullPolicy: &badNullPolicy,
	}
	err := q.Collapse(params)
	if err == nil {
		t.Fatal("expected error but got nothing")
	}
}
func TestCollapseValidParams(t *testing.T) {
	q := NewQuery(nil)
	np := NullPolicyIgnore
	hint := HintTopFC
	params := &CollapseParams{
		Field:      "field",
		Sort:       "field asc",
		NullPolicy: &np,
		Hint:       &hint,
	}
	err := q.Collapse(params)
	if err != nil {
		t.Fatalf("should be fine but got error %s", err)
	}
}

func TestExpand(t *testing.T) {
	q := NewQuery(nil)
	opts := &ExpandOptions{
		Sort: "field asc",
	}
	q.Expand(opts)
	if q.params.Get("expand") == "" {
		t.Fatal("expand param not registered")
	}
	if q.params.Get("expand.sort") == "" {
		t.Fatal("expand.sort param not registered")
	}
}

func TestSetQueryFields(t *testing.T) {
	q := NewQuery(nil)
	q.SetQueryFields([]string{"key, value"})
	if q.params.Get("qf") == "" {
		t.Fatal("qf param not registered")
	}
}

func TestSetMinimumShouldMatch(t *testing.T) {
	q := NewQuery(nil)
	q.SetMinimumShouldMatch("value")
	if q.params.Get("mm") == "" {
		t.Fatal("mm param not registered")
	}
}

func TestSetBoostFunctions(t *testing.T) {
	q := NewQuery(nil)
	q.SetBoostFunctions("abs(key)")
	if q.params.Get("bf") == "" {
		t.Fatal("bf param not registered")
	}
}

func TestSetBoostQuery(t *testing.T) {
	q := NewQuery(nil)
	q.SetBoostQuery("key:value")
	if q.params.Get("bq") == "" {
		t.Fatal("bq param not registered")
	}
}

func TestSetBoost(t *testing.T) {
	q := NewQuery(nil)
	q.SetBoost("mul(a,b)")
	if q.params.Get("boost") == "" {
		t.Fatal("boost param not registered")
	}
}

func TestSetUserFields(t *testing.T) {
	q := NewQuery(nil)
	q.SetUserFields([]string{"field1", "field2"})
	if q.params.Get("uf") == "" {
		t.Fatal("uf param not registered")
	}
}

func TestAddFacet(t *testing.T) {
	q := NewQuery(nil)
	f := &Facet{
		Field:        "field",
		Prefix:       "v",
		Contains:     "v",
		Limit:        10,
		MinCount:     5,
		Missing:      false,
		ExcludeTerms: []string{"term1", "term2"},
	}
	q.AddFacet(f)
	if q.params.Get("facet") == "" {
		t.Fatal("facet param not registered")
	}
	if q.params.Get("facet.field") == "" {
		t.Fatal("facet.field param not registered")
	}
	if q.params.Get("f.field.facet.limit") == "" {
		t.Fatal("f.field.facet.limit param not registered")
	}
	if q.params.Get("f.field.facet.mincount") == "" {
		t.Fatal("f.field.facet.mincount param not registered")
	}
	if q.params.Get("f.field.facet.prefix") == "" {
		t.Fatal("f.field.facet.prefix param not registered")
	}
	if q.params.Get("f.field.facet.contains") == "" {
		t.Fatal("f.field.facet.contains param not registered")
	}
	if q.params.Get("f.field.facet.excludeTerms") == "" {
		t.Fatal("f.field.facet.excludeTerms param not registered")
	}
}

func TestAddFacetPivot(t *testing.T) {
	q := NewQuery(nil)
	q.AddFacetPivot("field1.field2", 20)
	if q.params.Get("facet") == "" {
		t.Fatal("facet param not registered")
	}
	if q.params.Get("facet.pivot") == "" {
		t.Fatal("facet.pivot param not registered")
	}
	if q.params.Get("facet.pivot.mincount") == "" {
		t.Fatal("facet.pivot.mincount param not registered")
	}
}

func TestGroupNoParams(t *testing.T) {
	q := NewQuery(nil)
	err := q.Group(nil)
	if err == nil {
		t.Fatalf("expected error but got nothing")
	}
}

func TestGroupNoField(t *testing.T) {
	q := NewQuery(nil)
	params := &GroupParams{
		Field: "",
	}
	err := q.Group(params)
	if err == nil {
		t.Fatalf("expected error but got nothing")
	}
}

func TestGroupValidField(t *testing.T) {
	q := NewQuery(nil)
	params := &GroupParams{
		Field:            "field",
		Limit:            10,
		Offset:           5,
		Sort:             "field asc",
		ShowGroupsNumber: true,
	}
	err := q.Group(params)
	if err != nil {
		t.Fatalf("expected no error but got %s", err)
	}
	if q.params.Get("group") == "" {
		t.Fatal("group param not registered")
	}
	if q.params.Get("group.field") == "" {
		t.Fatal("group.field param not registered")
	}
	if q.params.Get("group.limit") == "" {
		t.Fatal("group.limit param not registered")
	}
	if q.params.Get("group.offset") == "" {
		t.Fatal("group.offset param not registered")
	}
	if q.params.Get("group.sort") == "" {
		t.Fatal("group.sort param not registered")
	}
	if q.params.Get("group.ngroups") == "" {
		t.Fatal("group.ngroups param not registered")
	}
}

func TestGroupValidQuery(t *testing.T) {
	q := NewQuery(nil)
	params := &GroupParams{
		Query: []string{"filed:val", "field:!val"},
	}
	err := q.Group(params)
	if err != nil {
		t.Fatalf("expected no error but got %s", err)
	}
	if q.params.Get("group") == "" {
		t.Fatal("group param not registered")
	}
	if q.params.Get("group.query") == "" {
		t.Fatal("group.query param not registered")
	}
}

func TestGroupValidFunc(t *testing.T) {
	q := NewQuery(nil)
	params := &GroupParams{
		Func: []string{"func1", "func2"},
	}
	err := q.Group(params)
	if err != nil {
		t.Fatalf("expected no error but got %s", err)
	}
	if q.params.Get("group") == "" {
		t.Fatal("group param not registered")
	}
	if q.params.Get("group.func") == "" {
		t.Fatal("group.func param not registered")
	}
}
