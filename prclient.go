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
	primary     *Connection
	replica     *Connection
	BasePath    string
	PrimaryPath string
	ReplicaPath string
}

// NewPrimaryReplicaClient returns two connections from the provided host and cores, one for the primary
// server and another for the replica. By default it is assumed that the primary server is used for
// writing data, and the replica server for reading data.
func NewPrimaryReplicaClient(ctx context.Context, pHost, pCore, rHost, rCore string, pClient, rClient *http.Client) (Client, error) {
	if pHost == "" || pCore == "" || rHost == "" || rCore == "" {
		return nil, ErrInvalidConfig
	}
	pConn := &Connection{
		Host:       pHost,
		Core:       pCore,
		httpClient: pClient,
	}
	pBasePath := formatBasePath(pHost, pCore)
	rConn := &Connection{
		Host:       rHost,
		Core:       rCore,
		httpClient: rClient,
	}
	rBasePath := formatBasePath(rHost, rCore)
	solrClient := &PRClient{primary: pConn, replica: rConn, PrimaryPath: pBasePath, ReplicaPath: rBasePath}
	return solrClient, nil
}

// SetBasicAuth sets auth credentials if needed.
func (c *PRClient) SetBasicAuth(username, password string) {
	c.primary.Username = username
	c.replica.Username = username
	c.primary.Password = password
	c.replica.Password = password
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
func (c *PRClient) Get(ctx context.Context, id string) (*Response, error) {
	vals := make(url.Values)
	vals.Set("id", id)
	url := c.formatReplicaURL("/get", vals.Encode())
	return read(ctx, c.replica, url)
}

// BatchGet ...
func (c *PRClient) BatchGet(ctx context.Context, ids []string, filter string) (*Response, error) {
	vals := make(url.Values)
	vals.Set("ids", strings.Join(ids, ","))
	vals.Set("fq", filter)
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
func (c *PRClient) CustomUpdate(ctx context.Context, item *UpdateBuilder) (*Response, error) {
	url := c.formatPrimaryURL("/update", "")
	return customUpdate(ctx, c.primary, url, item)
}
