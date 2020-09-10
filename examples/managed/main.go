package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/mecenat/solr"
)

func main() {
	ctx := context.Background()

	// To test the managed resources we need a supporting field type. Let's say for example
	// we want to have a managed synonyms list. We create the equivalent field type in our
	// core schema.
	sa, err := solr.NewSchemaAPI(ctx, "http://localhost:8983", "films", http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}

	sourcesField := &solr.FieldType{
		Name:  "source",
		CLass: "solr.TextField",
		IndexAnalyzer: &solr.Analyzer{
			Tokenizer: map[string]interface{}{
				"class": "solr.StandardTokenizerFactory",
			},
		},
		QueryAnalyzer: &solr.Analyzer{
			Tokenizer: map[string]interface{}{
				"class": "solr.StandardTokenizerFactory",
			},
			Filters: []map[string]interface{}{
				{
					"class":   "solr.ManagedSynonymGraphFilterFactory",
					"managed": "sources",
				},
				{
					"class": "solr.LowerCaseFilterFactory",
				},
			},
		},
	}

	scRes, err := sa.ReplaceFieldType(ctx, sourcesField)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(scRes.Header)

	// Initialize a new solr Managed Admin API
	ma, err := solr.NewManagedAPI(ctx, "http://localhost:8983", "films", http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}

	res, err := ma.RestManager(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Resources)

	res, err = ma.SynonymAddOptimal(ctx, "sources", []string{"TV", "television", "tele"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	res, err = ma.SynonymAdd(ctx, "sources", map[string][]string{"disc": {"dvd", "bluray"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	ign := map[string]interface{}{"ignoreCase": true}

	res, err = ma.SetInitArgs(ctx, "/analysis/synonyms/sources", ign)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	res, err = ma.SynonymGet(ctx, "sources", "disc")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	res, err = ma.SynonymDelete(ctx, "sources", "disc")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	res, err = ma.SynonymList(ctx, "sources")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Synonyms.ManagedMap)

	// in order for our edits to be saved we need to reload the core, using the CoreAPI
	ca, err := solr.NewCoreAdmin(ctx, "http://localhost:8983", http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}

	caRes, err := ca.Reload(ctx, "films")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(caRes.Header)
}
