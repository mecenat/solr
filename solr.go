// Copyright 2020 Mecenat (Authors: Konstantinos Koukouvis). All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package solr provides a solr client that enables the user to easily connect to
// one or more solr servers with support for the the basic CRUDL functionality
package solr

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// ErrInvalidConfig is returned when the hostname or corename are empty
var ErrInvalidConfig = errors.New("invalid configuration: no host or core provided")

// Connection represents the connection to the solr server and
// includes information about the address of the server and
// and the client to be used for connecting to it.
type Connection struct {
	httpClient *http.Client
	Host       string
	Core       string
}

// Client is the interface encompasing all the solr service methods
type Client interface {

	// Ping checks the connectivity of the solr server. It usually just returns with
	// Status = OK and a default response header, therefore this function just
	// returns an error in case there is no response, or an unexpected one
	Ping(ctx context.Context) error

	// Search performs a query to the solr server by using the `/select` endpoint, with the provided query
	// parameters. The query input can be easily created utilizing the provided helpers (check examples).
	// Currently only simple searches are supported.
	// For more info:
	// https://lucene.apache.org/solr/guide/8_5/overview-of-searching-in-solr.html
	Search(ctx context.Context, q *Query) (*Response, error)

	// Get performs a realtime get call to the solr server that returns the latest version of the document specified
	// by its id (uniqueKey field) without the associated cost of reopening a searcher. This is primarily useful
	// when using Solr as a NoSQL data store and not just a search index. For more info:
	// https://lucene.apache.org/solr/guide/8_5/realtime-get.html
	Get(ctx context.Context, id string) (*Response, error)

	// BatchGet performs a realtime get call to the solr server that returns the latest version of multiple documents
	// specified by their id (uniqueKey field) and filtered by the provided filter. The provided filter should
	// follow the format of the `fq` parameter but be concatenated in one string. For more info:
	// https://lucene.apache.org/solr/guide/8_5/realtime-get.html
	BatchGet(ctx context.Context, ids []string, filter string) (*Response, error)

	// Create adds a single document via JSON to the solr service. It calls the `/update/json/docs` endpoint.
	// Therefore the provided interface (item) must be a valid JSON object. This method accepts extra
	// options that are passed to the service as part of the request query. For more info:
	// https://lucene.apache.org/solr/guide/8_5/uploading-data-with-index-handlers.html#adding-a-single-json-document
	Create(ctx context.Context, item interface{}, opts *WriteOptions) (*Response, error)

	// BatchCreate adds multiple documents at once via JSON to the solr service. It calls the `/update` endpoint.
	// Therefore the provided interface (items) must be a valid array of JSON objects. This method accepts
	// extra options that are passed to the service as part of the request query. For more info:
	// https://lucene.apache.org/solr/guide/8_5/uploading-data-with-index-handlers.html#adding-multiple-json-documents
	BatchCreate(ctx context.Context, items interface{}, opts *WriteOptions) (*Response, error)

	// Update allows for partial updates of documents utilizing the "atomic" and the "in-place" updates approach.
	// The expected Fields input can be easily created using the provided helpers (check examples). This method
	// accepts extra options that are passed to the service as part of the request query. For more info:
	// https://lucene.apache.org/solr/guide/8_5/updating-parts-of-documents.html#atomic-updates
	Update(ctx context.Context, item *UpdatedFields, opts *WriteOptions) (*Response, error)

	// DeleteByID sends a JSON update command that deletes the document specified by its id (uniqueKey field).
	// It calls the `/update` endpoint and sends Solr JSON. This method accepts extra options that are
	// passed to the service as part of the request query. For more info:
	// https://lucene.apache.org/solr/guide/8_5/uploading-data-with-index-handlers.html#sending-json-update-commands
	DeleteByID(ctx context.Context, id string, opts *WriteOptions) (*Response, error)

	// DeleteByID sends a JSON update command that deletes the documents matching the given query. The query format
	// should follow the syntax of the Q parameter for the Search endpoint. It calls the `/update` endpoint and
	// sends Solr JSON. This method accepts extra options that are passed to the service as part of the
	// request query. For more info:
	// https://lucene.apache.org/solr/guide/8_5/uploading-data-with-index-handlers.html#sending-json-update-commands
	DeleteByQuery(ctx context.Context, query string, opts *WriteOptions) (*Response, error)

	// Clear is a helper method that removes all documents from the solr server. Use with caution.
	// It sends a DeleteByQuery request where the query is `*:*` and commit=true.
	Clear(ctx context.Context) (*Response, error)

	// Commit sends a JSON update command that commits all uncommited changes. Unless specified from one of the
	// options all write methods of this library will not commit their changes, therefore this method should
	// be called at the end of the a transaction to ensure that the indexes are properly updated.
	// For more info:
	// https://lucene.apache.org/solr/guide/8_5/uploading-data-with-index-handlers.html#sending-json-update-commands
	Commit(ctx context.Context, opts *CommitOptions) (*Response, error)

	// Rollback sends a JSON update command that rollbacks all uncommited changes. Unless specified from one of the
	// options all write methods of this library will not commit their changes, therefore this method should
	// be called if some action of the transaction returns an error and data cleaning is necessary.
	// For more info:
	// https://lucene.apache.org/solr/guide/8_5/uploading-data-with-index-handlers.html#sending-json-update-commands
	Rollback(ctx context.Context) (*Response, error)

	// Optimize sends a JSON update command that requests Solr to merge internal data structures. For a large index,
	// optimization will take some time to complete, but by merging many small segment files into larger segments, \
	// search performance may improve. More info:
	// https://lucene.apache.org/solr/guide/8_5/uploading-data-with-index-handlers.html#commit-and-optimize-during-updates
	Optimize(ctx context.Context, opts *OptimizeOptions) (*Response, error)

	// CustomUpdate allows the creation of a request to the `/update` endpoint that can include more than one update
	// command or for those that want a more finegrained request.
	CustomUpdate(ctx context.Context, item *UpdateBuilder) (*Response, error)
}

func read(ctx context.Context, client *http.Client, url string) (*Response, error) {
	return request(ctx, client, http.MethodGet, url, nil)
}

func create(ctx context.Context, client *http.Client, url string, item interface{}) (*Response, error) {
	bodyBytes, err := interfaceToBytes(item)
	if err != nil {
		return nil, err
	}

	err = isJSON(bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("Invalid JSON provided: %s", err)
	}

	return request(ctx, client, http.MethodPost, url, bodyBytes)
}

func batchCreate(ctx context.Context, client *http.Client, url string, items interface{}) (*Response, error) {
	bodyBytes, err := interfaceToBytes(items)
	if err != nil {
		return nil, err
	}

	err = isArrayOfJSON(bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("Invalid Array of JSON provided: %s", err)
	}

	return request(ctx, client, http.MethodPost, url, bodyBytes)
}

func update(ctx context.Context, client *http.Client, url string, item *UpdatedFields) (*Response, error) {
	ub := NewUpdateBuilder()
	ub.Add(item.fields)

	bodyBytes, err := interfaceToBytes(ub.commands)
	if err != nil {
		return nil, err
	}

	return request(ctx, client, http.MethodPost, url, bodyBytes)
}

func delete(ctx context.Context, client *http.Client, url string, doc Doc) (*Response, error) {
	ub := NewUpdateBuilder()
	ub.Delete(doc)

	bodyBytes, err := interfaceToBytes(ub.commands)
	if err != nil {
		return nil, err
	}

	return request(ctx, client, http.MethodPost, url, bodyBytes)
}

func commit(ctx context.Context, client *http.Client, url string, opts *CommitOptions) (*Response, error) {
	ub := NewUpdateBuilder()
	ub.Commit(opts)

	bodyBytes, err := interfaceToBytes(ub.commands)
	if err != nil {
		return nil, err
	}

	return request(ctx, client, http.MethodPost, url, bodyBytes)
}

func optimize(ctx context.Context, client *http.Client, url string, opts *OptimizeOptions) (*Response, error) {
	ub := NewUpdateBuilder()
	ub.Optimize(opts)

	bodyBytes, err := interfaceToBytes(ub.commands)
	if err != nil {
		return nil, err
	}

	return request(ctx, client, http.MethodPost, url, bodyBytes)
}

func rollback(ctx context.Context, client *http.Client, url string) (*Response, error) {
	ub := NewUpdateBuilder()
	ub.Rollback()

	bodyBytes, err := interfaceToBytes(ub.commands)
	if err != nil {
		return nil, err
	}

	return request(ctx, client, http.MethodPost, url, bodyBytes)
}

func customUpdate(ctx context.Context, client *http.Client, url string, item *UpdateBuilder) (*Response, error) {
	bodyBytes, err := interfaceToBytes(item.commands)
	if err != nil {
		return nil, err
	}

	return request(ctx, client, http.MethodPost, url, bodyBytes)
}
