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
	// Initialize a new solr schema API
	sa, err := solr.NewSchemaAPI(ctx, "http://localhost:8983", "films", http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve the schema
	res, err := sa.RetrieveSchema(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Header)
	// the below example shows how to get search through the different
	// fieldtypes in the schema and calculate the total number of
	// types that enable indexing by default.
	typesThatEnableIndexing := 0
	for _, f := range res.Schema.FieldTypes {
		if f.Name == "text_general" {
			fmt.Println(f.Name, readBoolPointer(f.MultiValued), readBoolPointer(f.OmitNorms))
		}
		if readBoolPointer(f.Indexed) {
			typesThatEnableIndexing++
		}
	}
	fmt.Println(typesThatEnableIndexing)

	// Create a custom field type with one analyzer
	singleAnalyzerFilterType := &solr.FieldType{
		Name:                 "custom",
		CLass:                "solr.TextField",
		PositionIncrementGap: "100",
		Analyzer: &solr.Analyzer{
			Tokenizer: map[string]interface{}{
				"class": "solr.StandardTokenizerFactory",
			},
		},
		FieldDefaultProperties: solr.FieldDefaultProperties{
			OmitNorms: newTrue(),
			Indexed:   newFalse(),
		},
	}

	res, err = sa.AddFieldType(ctx, singleAnalyzerFilterType)
	if err != nil {
		// in case there is an error with one of the commands we gave
		// to the schema API we can check the error details for the
		// actual entity that caused the error
		if len(res.Error.Details) > 0 {
			for _, d := range res.Error.Details {
				fmt.Println(d.MoreInfo())
			}
		}
		// without checking we just get the errors and the command that caused them
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// Retrieve the field type that we just created
	ft, err := sa.GetFieldType(ctx, "custom")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ft.Name, readBoolPointer(ft.Indexed))

	// we decided to replace the current version of our custom fieldtype with a new
	// to turn indexing on by default (when replacing a fieldtype we need to send
	// all the attributes and not only the ones we want to change).
	singleAnalyzerFilterTypeV2 := &solr.FieldType{
		Name:                 "custom",
		CLass:                "solr.TextField",
		PositionIncrementGap: "100",
		Analyzer: &solr.Analyzer{
			Tokenizer: map[string]interface{}{
				"class": "solr.StandardTokenizerFactory",
			},
		},
		FieldDefaultProperties: solr.FieldDefaultProperties{
			OmitNorms: newTrue(),
			Indexed:   newTrue(),
		},
	}

	res, err = sa.ReplaceFieldType(ctx, singleAnalyzerFilterTypeV2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// Retrieve the field type that we just created
	ft, err = sa.GetFieldType(ctx, "custom")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ft.Name, readBoolPointer(ft.Indexed))

	// Delete the field type that we just created
	res, err = sa.DeleteFieldType(ctx, "custom")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// the example below shows the creation of a fieldtype with two
	// analyzers, one for indexing and one for querying
	dualAnalyzerFilterType := &solr.FieldType{
		Name:  "nametext",
		CLass: "solr.TextField",
		IndexAnalyzer: &solr.Analyzer{
			Tokenizer: map[string]interface{}{
				"class": "solr.StandardTokenizerFactory",
			},
			Filters: []map[string]interface{}{
				{
					"class": "solr.LowerCaseFilterFactory",
				},
				{
					"class": "solr.KeepWordFilterFactory",
					"words": "protwords.txt",
				},
				{
					"class":    "solr.SynonymFilterFactory",
					"synonyms": "synonyms.txt",
				},
			},
		},
		QueryAnalyzer: &solr.Analyzer{
			Tokenizer: map[string]interface{}{
				"class": "solr.StandardTokenizerFactory",
			},
			Filters: []map[string]interface{}{
				{
					"class": "solr.LowerCaseFilterFactory",
				},
			},
		},
	}

	res, err = sa.AddFieldType(ctx, dualAnalyzerFilterType)
	if err != nil {
		if len(res.Error.Details) > 0 {
			for _, d := range res.Error.Details {
				fmt.Println(d.MoreInfo())
			}
		}
		log.Fatal(err)
	}

	// delete the field type we just created
	res, err = sa.DeleteFieldType(ctx, "nametext")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

}

// Helper functions for *bool handling

func readBoolPointer(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}

func newTrue() *bool {
	b := true
	return &b
}

func newFalse() *bool {
	b := false
	return &b
}
