package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mecenat/solr"
	"github.com/mecenat/solr/examples/data"
	"golang.org/x/net/context"
)

// Usually updates to the solr index will happen one by one. In rare occasions thought there might
// be a need of accomplishing more than one action at the same time.
// The following example showcases the use of a custom update. With that method it is possible
// to send more than one update command to the solr server, using the UpdateBuilder.
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

	// Documents to be used in the update request body
	filmDoc, err := data.Films[2].ToMap()
	if err != nil {
		log.Fatal(err)
	}

	delDoc := map[string]interface{}{
		"id": "3",
	}
	delDoc2 := map[string]interface{}{
		"query": "genre:horror",
	}
	delDoc3 := map[string]interface{}{
		"query": "*:*",
	}

	// Create a new update builder. The builder provides helpful
	// methods to easily create a complex json update body for
	// the a /update request
	ub := solr.NewUpdateBuilder()

	// Delete all the documents
	ub.Delete(delDoc3)
	// Rollback
	ub.Rollback()
	// Delete document with id "3"
	ub.Delete(delDoc)
	// Add the document with id "3"
	ub.Add(filmDoc)
	// Delete all the documents which contain "horror" as one of their genres
	ub.Delete(delDoc2)
	// Commit everything, with the option of expunging deletes
	ub.Commit(&solr.CommitOptions{ExpungeDeletes: true})
	// Optimize the index
	ub.Optimize(nil)

	// Send the custom update request
	res, err = slr.CustomUpdate(ctx, ub)
	if err != nil {
		log.Fatal(err)
	}

	// Clear the database, playtime is over
	res, err = slr.Clear(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
