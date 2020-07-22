package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mecenat/solr"
	"github.com/mecenat/solr/examples/data"
)

func main() {
	ctx := context.Background()
	slr, err := solr.NewSingleClient(ctx, "http://localhost:8983", "films", http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}

	// We start by inserting all the documents in order to have something
	// to play with
	res, err := slr.BatchCreate(ctx, data.Films, &solr.WriteOptions{Commit: true})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// Create a new query using an options object where we ask for a row limit
	// of 10 and enable Query type debugging.
	opts := &solr.ReadOptions{Rows: 10, Debug: solr.DebugTypeQuery}
	q := solr.NewQuery(opts)
	// We set the Q param to return everything
	q.SetQuery("*:*")
	// But filter on any film of the horror genre
	q.AddFilter("genre", "horror")
	// Then we set the sorting to happen descending based on the year property
	q.SetSort("year desc")

	fmt.Println(q.String())

	// We fire a search providing as input our Query
	res, err = slr.Search(ctx, q)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Data.NumFound)
	fmt.Println(len(res.Data.Docs))

	// the original query had a limit of 10 rows, therefore we
	// can ask solr for the remainder rows by setting the start
	// param of the already existing query
	q.SetStart(10)

	fmt.Println(q.String())

	// and fire a new search
	res, err = slr.Search(ctx, q)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Data.NumFound)
	fmt.Println(len(res.Data.Docs))

	// -----------

	// Create a new query without any options
	q2 := solr.NewQuery(nil)
	// We set the Q param to return only the horror films
	q2.SetQuery("genre:horror")
	// and then filter on those that are also a comedy
	q2.AddFilter("genre", "comedy")
	// we only want to return the name of the films
	q2.AddField("name")

	fmt.Println(q2.String())

	// Fire a search with that Query
	res, err = slr.Search(ctx, q2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Data.NumFound)
	fmt.Println(len(res.Data.Docs))

	// The Docs & Doc structs offer a helper ToBytes method
	// which can easily help you unmarshal them to your
	// structs
	var films []*data.Film

	fBytes, err := res.Data.Docs.ToBytes()
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(fBytes, &films)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(films[0].Name)

	// -----------

	// Create a new query without any options
	q3 := solr.NewQuery(nil)
	// We can search for movies that are both horror
	// and adventures using the following syntax
	q3.AddQuery("genre", "horror")
	q3.AddQuery("genre", "adventure")
	// we must set the operation here to AND since
	// the default in Solr is OR
	q3.SetOperationAND()

	fmt.Println(q3.String())

	// Fire a search with that Query
	res, err = slr.Search(ctx, q3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Data.NumFound)

	// -----------

	// The above will essentially create a query of the type:
	// q=genre:horror, genre:adventure&q.op=AND
	// if we want to do more advanced searches then we
	// should use SetQuery instead like so:
	// q=genre:horror AND (genre:comedy OR genre:action)

	q4 := solr.NewQuery(nil)
	q4.SetQuery("genre:horror AND (genre:comedy OR genre:action)")

	fmt.Println(q4.String())

	// Fire a search with that Query
	res, err = slr.Search(ctx, q4)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Data.NumFound)

	// -----------

	q5 := solr.NewQuery(&solr.ReadOptions{
		DefType: solr.DefTypeEDisMax,
	})
	q5.SetQuery("*:*")
	// Get all fields as well as score
	q5.AddField("*,score")
	q5.Collapse(&solr.CollapseParams{
		Field: "year",
	})
	q5.Expand(nil)

	fmt.Println(q5.String())

	// Fire a search with that Query
	res, err = slr.Search(ctx, q5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Expanded["1979"].NumFound)

	// -----------

	q6 := solr.NewQuery(&solr.ReadOptions{
		Rows: 2,
	})
	q6.SetQuery("*:*")

	q6.Group(&solr.GroupParams{
		Field:            "year",
		ShowGroupsNumber: true,
	})

	fmt.Println(q6.String())

	// Fire a search with that Query
	res, err = slr.Search(ctx, q6)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Grouped.ByField["year"].Matches)
	fmt.Println(res.Grouped.ByField["year"].NumberOfGroups)
	i := res.Grouped.ByField["year"].Groups[0].Value.(string)
	no := res.Grouped.ByField["year"].Groups[0].DocList.NumFound
	fmt.Println(i, no)

	// -----------

	q7 := solr.NewQuery(&solr.ReadOptions{
		Rows: 1,
	})
	q7.SetQuery("*:*")

	q7.AddFacet(&solr.Facet{
		Field:    "genre",
		MinCount: 1,
	})

	q7.AddFacetPivot("genre,directed_by", 1)

	fmt.Println(q7.String())

	// Fire a search with that Query
	res, err = slr.Search(ctx, q7)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("horror", res.FacetCounts.Fields.Get("genre")["horror"])
	fmt.Println(res.FacetCounts.Pivot["genre,directed_by"][0].Pivot[0].Count)

	// -----------

	q8 := solr.NewQuery(&solr.ReadOptions{
		Rows: 2,
	})
	q8.SetQuery("*:*")

	q8.Group(&solr.GroupParams{
		Query: []string{"year:1968", "year:(!1968)"},
	})

	fmt.Println(q8.String())

	// Fire a search with that Query
	res, err = slr.Search(ctx, q8)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Grouped)
	no = res.Grouped.ByQueryOrFunc["year:(!1968)"].DocList.NumFound
	fmt.Println("!year:1968", no)

	// Clear the database, playtime is over
	res, err = slr.Clear(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
