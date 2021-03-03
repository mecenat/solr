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
	conn, err := solr.NewConnection("http://localhost:8983", "films", http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}

	slr, err := solr.NewSingleClient(conn)
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

	// Create a new update builder. The builder provides helpful
	// methods to easily create a complex json update body for
	// the a /update request
	ub := solr.NewUpdateBuilder()

	// Add 4 films (will overwrite by default)
	ub.Add(data.Films[2])
	ub.Add(data.Films[3])
	ub.Add(data.Films[4])
	ub.Add(data.Films[5])

	// Delete films with IDs 3 & 4 but also all those that are horror films.
	ub.DeleteByID(data.Films[3].ID)
	ub.DeleteByID(data.Films[4].ID)
	ub.DeleteByQuery("genre:horror")

	// Change the name of 3 of the films and update it
	data.Films[2].Name = "New Name"
	data.Films[6].Name = "New Name"
	data.Films[8].Name = "New Name"
	ub.Add(data.Films[2])
	ub.Add(data.Films[6])
	ub.Add(data.Films[8])

	// Send the custom update request
	res, err = slr.CustomUpdate(ctx, ub, &solr.WriteOptions{Commit: true})
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
