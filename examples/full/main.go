package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mecenat/solr"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	slr, err := solr.NewSingleClient("http://localhost:8983", "films", http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}

	// Create a document entry for the first film of the array
	res, err := slr.Create(ctx, films[0], &solr.WriteOptions{Commit: true})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// Update the genre attribute
	uf := solr.NewUpdateDocument("1")
	uf.Add("genre", "Slasher")
	uf.Add("directed_by", "some guy")
	res, err = slr.Update(ctx, uf, &solr.WriteOptions{Commit: true})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// Then delete that document
	res, err = slr.DeleteByID(ctx, "1", &solr.WriteOptions{Commit: true})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// Create a document entry for the first film of the array
	// but without commiting anything
	res, err = slr.Create(ctx, films[0], nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// RealTime Get the film that was just added
	res, err = slr.Get(ctx, films[0].ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Doc)

	// rollback the latest changes
	res, err = slr.Rollback(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)

	// RealTime Get the film that was just added
	res, err = slr.Get(ctx, films[0].ID)
	if err != nil {
		log.Fatal(err)
	}
	// should be null since we rollbacked
	fmt.Println(res.Doc)

	// insert all the documents at once
	res, err = slr.BatchCreate(ctx, films, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// RealTime Get the film that was just added
	res, err = slr.BatchGet(ctx, []string{films[0].ID, films[3].ID}, "")
	if err != nil {
		log.Fatal(err)
	}
	// should be 2
	fmt.Println(res.Data.NumFound)

	// commit the latest changes
	res, err = slr.Commit(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)

	// query for all items and filter only to horror films
	opts := &solr.ReadOptions{Rows: 10}
	q := solr.NewQuery(opts)
	q.SetQuery("*:*")
	q.AddFilter("genre", "horror")
	q.SetSort("year desc")
	fmt.Println(q.String())
	res, err = slr.Search(ctx, q)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Data.NumFound)
	fmt.Println(len(res.Data.Docs))

	// the original query had a limit to 10 rows, therefore we
	// can ask solr for the remainder rows by setting the start
	// param of the already existing query
	q.SetStart(10)
	fmt.Println(q.String())
	res, err = slr.Search(ctx, q)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Data.NumFound)
	fmt.Println(len(res.Data.Docs))

	// we can delete all horror films by specifying a query
	// for deletion
	res, err = slr.DeleteByQuery(ctx, "genre:horror", &solr.WriteOptions{Commit: true})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// we query for all items
	q2 := solr.NewQuery(opts)
	q2.SetQuery("*:*")
	res, err = slr.Search(ctx, q2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Data.NumFound)
	fmt.Println(len(res.Data.Docs))

	remainingFilms := []*film{}
	resDocBytes, err := res.Data.Docs.ToBytes()
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(resDocBytes, &remainingFilms)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(remainingFilms))
	fmt.Println(remainingFilms[0].Name)

	// the delete everything in films
	res, err = slr.Clear(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
