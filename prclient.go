package solr

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// PRClient implements the solr interface in Primary - Replica server
// architecture. It contains a connection to a Primary server used
// for writing data, and a connection to a Replica server used
// for reading data.
type PRClient struct {
	primary     connection
	replica     connection
	PrimaryPath string
	ReplicaPath string
}

// NewPrimaryReplicaClient returns two connections from the provided host and cores, one for the primary
// server and another for the replica. By default it is assumed that the primary server is used for
// writing data, and the replica server for reading data.
func NewPrimaryReplicaClient(primaryConn, replicaConn connection) (Client, error) {
	pBasePath := primaryConn.formatBasePath()
	rBasePath := replicaConn.formatBasePath()
	return &PRClient{
		primary:     primaryConn,
		replica:     replicaConn,
		PrimaryPath: pBasePath,
		ReplicaPath: rBasePath,
	}, nil
}

// SetBasicAuth sets auth credentials if needed.
func (c *PRClient) SetBasicAuth(username, password string) {
	c.primary.setBasicAuth(username, password)
	c.replica.setBasicAuth(username, password)
}

func (c *PRClient) formatPrimaryURL(path string, query string) string {
	if query != "" {
		return c.PrimaryPath + path + "?" + query
	}
	return c.PrimaryPath + path
}

func (c *PRClient) formatReplicaURL(path string, query string) string {
	if query != "" {
		return c.ReplicaPath + path + "?" + query
	}
	return c.ReplicaPath + path
}

// Ping tests the connectivity of both servers
func (c *PRClient) Ping(ctx context.Context) error {
	url := c.formatPrimaryURL("/admin/ping", "")
	res, err := c.primary.request(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if res.Status != nil && *res.Status != "OK" {
		return fmt.Errorf("error pinging primary server, status: %s", *res.Status)
	}
	url = c.formatReplicaURL("/admin/ping", "")
	res, err = c.replica.request(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if res.Status != nil && *res.Status != "OK" {
		return fmt.Errorf("error pinging replica server, status: %s", *res.Status)
	}
	return nil
}

// Search ...
func (c *PRClient) Search(ctx context.Context, q *Query) (*Response, error) {
	url := c.formatReplicaURL("/select", q.String())
	return read(ctx, c.replica, url)
}

// Get ...
func (c *PRClient) Get(ctx context.Context, id, filter string) (*Response, error) {
	vals := make(url.Values)
	vals.Set("id", id)
	if filter != "" {
		vals.Set("fq", filter)
	}
	url := c.formatReplicaURL("/get", vals.Encode())
	return read(ctx, c.replica, url)
}

// BatchGet ...
func (c *PRClient) BatchGet(ctx context.Context, ids []string, filter string) (*Response, error) {
	vals := make(url.Values)
	vals.Set("ids", strings.Join(ids, ","))
	if filter != "" {
		vals.Set("fq", filter)
	}
	url := c.formatReplicaURL("/get", vals.Encode())
	return read(ctx, c.replica, url)
}

// Create ...
func (c *PRClient) Create(ctx context.Context, item interface{}, opts *WriteOptions) (*Response, error) {
	url := c.formatPrimaryURL("/update/json/docs", opts.formatQueryFromOpts().Encode())
	return create(ctx, c.primary, url, item)
}

// BatchCreate ...
func (c *PRClient) BatchCreate(ctx context.Context, items interface{}, opts *WriteOptions) (*Response, error) {
	url := c.formatPrimaryURL("/update", opts.formatQueryFromOpts().Encode())
	return batchCreate(ctx, c.primary, url, items)
}

// Update ...
func (c *PRClient) Update(ctx context.Context, item *UpdatedFields, opts *WriteOptions) (*Response, error) {
	url := c.formatPrimaryURL("/update", opts.formatQueryFromOpts().Encode())
	return update(ctx, c.primary, url, item)
}

// Commit ...
func (c *PRClient) Commit(ctx context.Context, opts *CommitOptions) (*Response, error) {
	url := c.formatPrimaryURL("/update", "")
	return commit(ctx, c.primary, url, opts)
}

// Rollback ...
func (c *PRClient) Rollback(ctx context.Context) (*Response, error) {
	url := c.formatPrimaryURL("/update?commit=true", "")
	return rollback(ctx, c.primary, url)
}

// Optimize ...
func (c *PRClient) Optimize(ctx context.Context, opts *OptimizeOptions) (*Response, error) {
	url := c.formatPrimaryURL("/update", "")
	return optimize(ctx, c.primary, url, opts)
}

// DeleteByID ...
func (c *PRClient) DeleteByID(ctx context.Context, id string, opts *WriteOptions) (*Response, error) {
	url := c.formatPrimaryURL("/update", opts.formatQueryFromOpts().Encode())
	return delete(ctx, c.primary, url, formatDeleteByID(id))
}

// DeleteByQuery ...
func (c *PRClient) DeleteByQuery(ctx context.Context, query string, opts *WriteOptions) (*Response, error) {
	url := c.formatPrimaryURL("/update", opts.formatQueryFromOpts().Encode())
	return delete(ctx, c.primary, url, formatDeleteByQuery(query))
}

// Clear ...
func (c *PRClient) Clear(ctx context.Context) (*Response, error) {
	return c.DeleteByQuery(ctx, "*:*", &WriteOptions{Commit: true})
}

// CustomUpdate ...
func (c *PRClient) CustomUpdate(ctx context.Context, item *UpdateBuilder, opts *WriteOptions) (*Response, error) {
	url := c.formatPrimaryURL("/update", opts.formatQueryFromOpts().Encode())
	return customUpdate(ctx, c.primary, url, item)
}
