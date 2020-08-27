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
// host and core.
func NewSingleClient(ctx context.Context, host, core string, client *http.Client) (Client, error) {
	if host == "" || core == "" {
		return nil, ErrInvalidConfig
	}
	conn := &Connection{
		Host:       host,
		Core:       core,
		httpClient: client,
	}
	bp := formatBasePath(host, core)
	return &SingleClient{conn: conn, BasePath: bp}, nil
}

// SetAuth sets auth credentials if needed.
func (c *SingleClient) SetAuth(username, password string) {
	c.conn.Username = username
	c.conn.Password = password
}

func (c *SingleClient) formatURL(path string, query string) string {
	if query != "" {
		return c.BasePath + path + "?" + query
	}
	return c.BasePath + path
}

// Ping ...
func (c *SingleClient) Ping(ctx context.Context) error {
	url := c.formatURL("/admin/ping", "")
	res, err := c.conn.request(ctx, http.MethodGet, url, nil)
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
	url := c.formatURL("/select", q.String())
	return read(ctx, c.conn, url)
}

// Get ...
func (c *SingleClient) Get(ctx context.Context, id string) (*Response, error) {
	vals := make(url.Values)
	vals.Set("id", id)
	url := c.formatURL("/get", vals.Encode())
	return read(ctx, c.conn, url)
}

// BatchGet ...
func (c *SingleClient) BatchGet(ctx context.Context, ids []string, filter string) (*Response, error) {
	vals := make(url.Values)
	vals.Set("ids", strings.Join(ids, ","))
	vals.Set("fq", filter)
	url := c.formatURL("/get", vals.Encode())
	return read(ctx, c.conn, url)
}

// Create ...
func (c *SingleClient) Create(ctx context.Context, item interface{}, opts *WriteOptions) (*Response, error) {
	url := c.formatURL("/update/json/docs", opts.formatQueryFromOpts().Encode())
	return create(ctx, c.conn, url, item)
}

// BatchCreate ...
func (c *SingleClient) BatchCreate(ctx context.Context, items interface{}, opts *WriteOptions) (*Response, error) {
	url := c.formatURL("/update", opts.formatQueryFromOpts().Encode())
	return batchCreate(ctx, c.conn, url, items)
}

// Update ...
func (c *SingleClient) Update(ctx context.Context, item *UpdatedFields, opts *WriteOptions) (*Response, error) {
	url := c.formatURL("/update", opts.formatQueryFromOpts().Encode())
	return update(ctx, c.conn, url, item)
}

// Commit ...
func (c *SingleClient) Commit(ctx context.Context, opts *CommitOptions) (*Response, error) {
	url := c.BasePath + "/update"
	return commit(ctx, c.conn, url, opts)
}

// Rollback ...
func (c *SingleClient) Rollback(ctx context.Context) (*Response, error) {
	url := c.BasePath + "/update?commit=true"
	return rollback(ctx, c.conn, url)
}

// Optimize ...
func (c *SingleClient) Optimize(ctx context.Context, opts *OptimizeOptions) (*Response, error) {
	url := c.formatURL("/update", "")
	return optimize(ctx, c.conn, url, opts)
}

// DeleteByID ...
func (c *SingleClient) DeleteByID(ctx context.Context, id string, opts *WriteOptions) (*Response, error) {
	url := c.formatURL("/update", opts.formatQueryFromOpts().Encode())
	return delete(ctx, c.conn, url, formatDeleteByID(id))
}

// DeleteByQuery ...
func (c *SingleClient) DeleteByQuery(ctx context.Context, query string, opts *WriteOptions) (*Response, error) {
	url := c.formatURL("/update", opts.formatQueryFromOpts().Encode())
	return delete(ctx, c.conn, url, formatDeleteByQuery(query))
}

// Clear ...
func (c *SingleClient) Clear(ctx context.Context) (*Response, error) {
	return c.DeleteByQuery(ctx, "*:*", &WriteOptions{Commit: true})
}

// CustomUpdate ...
func (c *SingleClient) CustomUpdate(ctx context.Context, item *UpdateBuilder) (*Response, error) {
	url := c.formatURL("/update", "")
	return customUpdate(ctx, c.conn, url, item)
}
