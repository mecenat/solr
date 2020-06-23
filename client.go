package solr

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SingleClient implements the solr interface and is the basic connection
// to a solr server.
type SingleClient struct {
	conn     *Connection
	BasePath string
}

// NewSingleClient returns a connection to the solr client provided by the given
// host and core. A ping is also sent to the server to verify that it is
// active and a connection can be made.
func NewSingleClient(host, core string, client *http.Client) (Client, error) {
	conn := &Connection{
		Host:       host,
		Core:       core,
		httpClient: client,
	}
	bp := formatBasePath(host, core)
	solrClient := &SingleClient{conn: conn, BasePath: bp}
	err := solrClient.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return solrClient, nil
}

func (c *SingleClient) formatURL(path string, query url.Values) string {
	if query != nil {
		return c.BasePath + path + "?" + query.Encode()
	}
	return c.BasePath + path
}

// Ping ...
func (c *SingleClient) Ping(ctx context.Context) error {
	url := c.formatURL("/admin/ping", nil)
	res, err := request(ctx, c.conn.httpClient, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if res.Status != nil && *res.Status != "OK" {
		return fmt.Errorf("error pinging solr, status: %s", *res.Status)
	}
	return nil
}

// Search ...
func (c *SingleClient) Search(ctx context.Context, q *Query) (*Response, error) {
	url := c.formatURL("/select", q.params)
	return read(ctx, c.conn.httpClient, url)
}

// Get ...
func (c *SingleClient) Get(ctx context.Context, id string) (*Response, error) {
	vals := make(url.Values)
	vals.Set("id", id)
	url := c.formatURL("/get", vals)
	return read(ctx, c.conn.httpClient, url)
}

// BatchGet ...
func (c *SingleClient) BatchGet(ctx context.Context, ids []string, filter string) (*Response, error) {
	vals := make(url.Values)
	vals.Set("ids", strings.Join(ids, ","))
	vals.Set("fq", filter)
	url := c.formatURL("/get", vals)
	return read(ctx, c.conn.httpClient, url)
}

// Create ...
func (c *SingleClient) Create(ctx context.Context, item interface{}, opts *WriteOptions) (*Response, error) {
	url := c.formatURL("/update/json/docs", opts.formatQueryFromOpts())
	return create(ctx, c.conn.httpClient, url, item)
}

// BatchCreate ...
func (c *SingleClient) BatchCreate(ctx context.Context, items interface{}, opts *WriteOptions) (*Response, error) {
	url := c.formatURL("/update", opts.formatQueryFromOpts())
	return batchCreate(ctx, c.conn.httpClient, url, items)
}

// Update ...
func (c *SingleClient) Update(ctx context.Context, item *Fields, opts *WriteOptions) (*Response, error) {
	url := c.formatURL("/update", opts.formatQueryFromOpts())
	return update(ctx, c.conn.httpClient, url, item)
}

// Commit ...
func (c *SingleClient) Commit(ctx context.Context) (*Response, error) {
	url := c.BasePath + "/update"
	return commit(ctx, c.conn.httpClient, url)
}

// Rollback ...
func (c *SingleClient) Rollback(ctx context.Context) (*Response, error) {
	url := c.BasePath + "/update?commit=true"
	return rollback(ctx, c.conn.httpClient, url)
}

// DeleteByID ...
func (c *SingleClient) DeleteByID(ctx context.Context, id string, opts *WriteOptions) (*Response, error) {
	url := c.formatURL("/update", opts.formatQueryFromOpts())
	return delete(ctx, c.conn.httpClient, url, formatDeleteByID(id))
}

// DeleteByQuery ...
func (c *SingleClient) DeleteByQuery(ctx context.Context, query string, opts *WriteOptions) (*Response, error) {
	url := c.formatURL("/update", opts.formatQueryFromOpts())
	return delete(ctx, c.conn.httpClient, url, formatDeleteByQuery(query))
}

// Clear ...
func (c *SingleClient) Clear(ctx context.Context) (*Response, error) {
	return c.DeleteByQuery(ctx, "*:*", &WriteOptions{Commit: true})
}
