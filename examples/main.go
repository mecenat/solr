package main

import (
	"fmt"
	"net/http"

	"github.com/mecenat/solr"
	"gopkg.in/square/go-jose.v2/json"
)

// XXX: just for testing purposes, not part of actual repo
func main() {
	slr := solr.New("http://localhost:8983", "films", http.DefaultClient)
	sts, err := slr.Ping()
	fmt.Println(sts, err)

	// opts := &solr.QueryOptions{}
	// q := solr.NewQuery(opts)
	// q.AddQuery("*:*")
	// q.SetQuery("alias:two object_type:topbanner")
	// q.AddQuery("id:666")
	// q.SetOperation("and")
	// q.AddFilter("published", "true")
	// q.AddFilter("deleted", "false")
	// q.SetStart(1)
	// q.SetSort("created_at desc")
	// fmt.Println(q.String())
	// res, err := slr.Search(q)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(res.Header)
	// fmt.Println(res.Error)
	// fmt.Println(res.Data)
	// fmt.Println(res.Debug)

	// fmt.Println(res2.Error)
	// fmt.Println(res.Response.Docs[0])
	// its := []*mkt.NativeMarketingItem{}
	// err = json.Unmarshal(by, &its)
	// fmt.Println(its[1].Alias)

	f := &film{
		ID:       "snatch_2000",
		Name:     "Snatch",
		Year:     "2000",
		Genre:    []string{"Crime", "Comedy"},
		Director: []string{"Guy Ritchie"},
	}

	b, err := json.Marshal(f)
	res3, err := slr.Create(b)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res3)

	res2, err := slr.Get("snatch_2000")
	if err != nil {
		fmt.Println(err)
		return
	}
	var f2 film
	fby, err := res2.Doc.ToBytes()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(fby))

	err = json.Unmarshal(fby, &f2)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(f2.Name)
}

type film struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Year     string   `json:"year"`
	Genre    []string `json:"genre"`
	Director []string `json:"directed_by"`
}
