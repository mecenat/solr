package solr

import (
	"context"
	"fmt"
	"net/http"
)

type SolrClient interface {
	Ping() error
	BatchGet(ctx context.Context, ids []string, filter string) (*Response, error)
	BatchCreate(ctx context.Context, items interface{}, opts *WriteOptions) (*Response, error)
	Search(ctx context.Context, q *Query) (*Response, error)
	Get(ctx context.Context, id string) (*Response, error)
	Create(ctx context.Context, item interface{}, opts *WriteOptions) (*Response, error)
	Update(ctx context.Context, item *Fields, opts *WriteOptions) (*Response, error)
	DeleteByID(ctx context.Context, id string, opts *WriteOptions) (*Response, error)
	DeleteByQuery(ctx context.Context, query string, opts *WriteOptions) (*Response, error)
	Clear(ctx context.Context) (*Response, error)
	Commit(ctx context.Context) (*Response, error)
	Rollback(ctx context.Context) (*Response, error)
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

func update(ctx context.Context, client *http.Client, url string, item *Fields) (*Response, error) {
	ub := NewUpdateBuilder()
	ub.Add(formatDocEntry(item.fields))

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

func commit(ctx context.Context, client *http.Client, url string) (*Response, error) {
	ub := NewUpdateBuilder()
	ub.Commit()

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
