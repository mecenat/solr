package main

import (
	"fmt"
	"net/http"

	"github.com/mecenat/solr"
)

// XXX: just for testing purposes, not part of actual repo
func main() {
	slr := solr.New("http://localhost:8983", "nativemarketing", http.DefaultClient)
	sts, err := slr.Ping()
	fmt.Println(sts, err)

	opts := &solr.QueryOptions{}
	q := solr.NewQuery(opts)
	// q.AddQuery("*:*")
	q.SetQuery("alias:two object_type:topbanner")
	// q.AddQuery("id:666")
	q.SetOperation("and")
	// q.AddFilter("published", "true")
	// q.AddFilter("deleted", "false")
	// q.SetStart(1)
	// q.SetSort("created_at desc")
	fmt.Println(q.String())
	res, err := slr.Search(q)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Header)
	fmt.Println(res.Error)
	fmt.Println(res.Data)
	fmt.Println(res.Debug)

	res2, err := slr.Get([]string{"666", "b0de01fe-9689-11ea-b762-acde48001122"}, "")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res2.Data)
	fmt.Println(res2.Error)
	// fmt.Println(res.Response.Docs[0])
	// by, err := res.Response.Docs.ToBytes()
	// its := []*mkt.NativeMarketingItem{}
	// err = json.Unmarshal(by, &its)
	// fmt.Println(its[1].Alias)
}
