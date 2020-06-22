package solr

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	conn       *Connection
	BasePath   string
	AutoCommit bool
}

type Connection struct {
	httpClient *http.Client
	Host       string
	Core       string
}

func New(host, core string, client *http.Client) SolrClient {
	conn := &Connection{
		Host:       host,
		Core:       core,
		httpClient: client,
	}
	bp := formatBasePath(host, core)
	return &Client{conn: conn, BasePath: bp}
}

func (c *Client) formatURL(path string, query url.Values) string {
	if query != nil {
		return c.BasePath + path + "?" + query.Encode()
	}
	return c.BasePath + path
}

func (c *Client) request(ctx context.Context, method, url string, body []byte) (*Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := c.conn.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	var r Response
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(b))

	// err = json.NewDecoder(res.Body).Decode(&r)
	// if err != nil {
	// 	return nil, err
	// }

	// if r.Error != nil {
	// 	return nil, r.Error
	// }

	return &r, nil
}

// Ping checks the connectivity of the solr server. It usually just returns with
// Status = OK and a default response header, therefore this function just
// returns an error in case there is no response, or an unexpected one
func (c *Client) Ping() error {
	url := c.formatURL("/admin/ping", nil)
	res, err := c.request(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if res.Status != nil && *res.Status != "OK" {
		return fmt.Errorf("error pinging solr, status: %s", *res.Status)
	}
	return nil
}

// Search ...
func (c *Client) Search(ctx context.Context, q *Query) (*Response, error) {
	url := c.formatURL("/select", q.params)
	return read(ctx, c.conn.httpClient, url)
}

// Get performs a real time get that returns the latest version of the specified document(s)
// https://lucene.apache.org/solr/guide/8_5/realtime-get.html
func (c *Client) Get(ctx context.Context, id string) (*Response, error) {
	vals := make(url.Values)
	vals.Set("id", id)
	url := c.formatURL("/get", vals)
	return read(ctx, c.conn.httpClient, url)
}

// BatchGet performs a real time get that returns the latest version of the specified document(s)
// https://lucene.apache.org/solr/guide/8_5/realtime-get.html
func (c *Client) BatchGet(ctx context.Context, ids []string, filter string) (*Response, error) {
	vals := make(url.Values)
	vals.Set("ids", strings.Join(ids, ","))
	vals.Set("fq", filter)
	url := c.formatURL("/get", vals)
	return read(ctx, c.conn.httpClient, url)
}

// Create inserts the given interface to the db. The provided interface must be valid JSON
// implementing the JSONMarshaler interface (see /examples)
func (c *Client) Create(ctx context.Context, item interface{}, opts *WriteOptions) (*Response, error) {
	url := c.formatURL("/update/json/docs", opts.formatQueryFromOpts())
	return create(ctx, c.conn.httpClient, url, item)
}

// BatchCreate inserts an array of JSON data to solr. The provided interface must be a valid array
// of JSON objects
func (c *Client) BatchCreate(ctx context.Context, items interface{}, opts *WriteOptions) (*Response, error) {
	url := c.formatURL("/update", opts.formatQueryFromOpts())
	return batchCreate(ctx, c.conn.httpClient, url, items)
}

// Update does an atomic update
func (c *Client) Update(ctx context.Context, item *Fields, opts *WriteOptions) (*Response, error) {
	url := c.formatURL("/update", opts.formatQueryFromOpts())
	return update(ctx, c.conn.httpClient, url, item)
}

// Commit makes a call to the update endpoint with the commit option set
// to true commiting all uncommited changes
func (c *Client) Commit(ctx context.Context) (*Response, error) {
	url := c.BasePath + "/update"
	return commit(ctx, c.conn.httpClient, url)
}

// Rollback deletes all uncommited changes
func (c *Client) Rollback(ctx context.Context) (*Response, error) {
	url := c.BasePath + "/update?commit=true"
	return rollback(ctx, c.conn.httpClient, url)
}

// DeleteByID deletes the document specified by its id (uniqueKey field)
func (c *Client) DeleteByID(ctx context.Context, id string, opts *WriteOptions) (*Response, error) {
	url := c.formatURL("/update", opts.formatQueryFromOpts())
	return delete(ctx, c.conn.httpClient, url, formatDeleteByID(id))
}

// DeleteByQuery deletes the document(s) that are returned by the query
func (c *Client) DeleteByQuery(ctx context.Context, query string, opts *WriteOptions) (*Response, error) {
	url := c.formatURL("/update", opts.formatQueryFromOpts())
	return delete(ctx, c.conn.httpClient, url, formatDeleteByQuery(query))
}

// Clear is a wrapper function for DeleteByQuery where the query is "*:*" which erases
// ALL documents from the database. Use at your own risk!
func (c *Client) Clear(ctx context.Context) (*Response, error) {
	return c.DeleteByQuery(ctx, "*:*", &WriteOptions{Commit: true})
}
