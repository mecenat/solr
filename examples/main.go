package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/mecenat/solr"
)

// Running this file will also create the necessary schema
// using the schema API

func main() {
	ctx := context.Background()
	// Initialize a new solr schema API
	sa, err := solr.NewSchemaAPI(ctx, "http://localhost:8983", "films", http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}

	id := &solr.Field{
		Name: "id",
		Type: "string",
	}

	res, err := sa.AddField(ctx, id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	name := &solr.Field{
		Name: "name",
		Type: "string",
	}
	res, err = sa.AddField(ctx, name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	year := &solr.Field{
		Name: "year",
		Type: "string",
	}
	res, err = sa.AddField(ctx, year)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	genre := &solr.Field{
		Name: "genre",
		Type: "text_general",
	}
	res, err = sa.AddField(ctx, genre)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	directedBy := &solr.Field{
		Name: "directed_by",
		Type: "text_general",
	}
	res, err = sa.AddField(ctx, directedBy)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	seen := &solr.Field{
		Name: "seen_counter",
		Type: "pint",
	}
	res, err = sa.AddField(ctx, seen)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)
}
