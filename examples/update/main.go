package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mecenat/solr"
	"github.com/mecenat/solr/examples/data"
	"golang.org/x/net/context"
)

// The following example acts as a showcase of most of the methods provided
// by this library, methods that allow the usage of the CRUD paradigm.

func main() {
	ctx := context.Background()
	slr, err := solr.NewSingleClient("http://localhost:8983", "data.films", http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}

	// Create a document entry for the first film of the array
	res, err := slr.Create(ctx, data.Films[0], &solr.WriteOptions{Commit: true})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// Update some of the attributes of the document with id = 1
	uf := solr.NewUpdateDocument("1")
	// uf.Add("genre", []string{"Crime", "Slasher"})
	// uf.AddDistinct("genre", []string{"Crime", "Slasher", "Commedy"})
	// uf.Remove("genre", "Commedy")
	// Avoid doing as the above, this will only sent the last action
	// use Set instead
	uf.Set("genre", []string{"Crime", "Action", "Comedy"})
	uf.Add("directed_by", "Some guy")
	uf.IncrementBy("seen_counter", 1)
	res, err := slr.Update(ctx, uf, &solr.WriteOptions{Commit: true})
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
	res, err = slr.Create(ctx, data.Films[0], nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// RealTime Get the film that was just added
	res, err = slr.Get(ctx, data.films[0].ID)
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
	res, err = slr.Get(ctx, data.films[0].ID)
	if err != nil {
		log.Fatal(err)
	}
	// should be null since we rollbacked
	fmt.Println(res.Doc)

	// insert all the documents at once
	res, err = slr.BatchCreate(ctx, data.films, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// RealTime Get the film that was just added
	res, err = slr.BatchGet(ctx, []string{data.films[0].ID, data.films[3].ID}, "")
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

	// we can delete all horror films by specifying a query
	// for deletion
	res, err = slr.DeleteByQuery(ctx, "genre:horror", &solr.WriteOptions{Commit: true})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// if we ever feel the need to delete all the documents
	// in the core, there's a helper for that!
	res, err = slr.Clear(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
