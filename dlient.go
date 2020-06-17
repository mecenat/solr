package solr

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	conn     *Connection
	BasePath string
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

func (c *Client) Ping() (*Response, error) {
	url := c.BasePath + "/admin/ping"
	return request(context.Background(), *c.conn, http.MethodGet, url, nil)
	// if err != nil {
	// 	return 0, err
	// }
	// return res.StatusCode, nil
}

func (c *Client) Search(q *Query) (*Response, error) {
	url := c.BasePath + "/select?" + q.String()
	return request(context.Background(), *c.conn, http.MethodGet, url, nil)
}

// Get performs a real time get that returns the latest version of the specified document(s)
// https://lucene.apache.org/solr/guide/8_5/realtime-get.html
func (c *Client) Get(id string) (*Response, error) {
	query := fmt.Sprintf("?id=%s", id)
	url := c.BasePath + "/get" + query
	return request(context.Background(), *c.conn, http.MethodGet, url, nil)
}

// BatchGet performs a real time get that returns the latest version of the specified document(s)
// https://lucene.apache.org/solr/guide/8_5/realtime-get.html
func (c *Client) BatchGet(ids []string, filter string) (*Response, error) {
	query := fmt.Sprintf("?ids=%s&fq=%s", strings.Join(ids, ","), url.QueryEscape(filter))
	url := c.BasePath + "/get" + query
	return request(context.Background(), *c.conn, http.MethodGet, url, nil)
}

func (c *Client) Create(item []byte) (*Response, error) {
	url := c.BasePath + "/update?json.command=false"
	return request(context.Background(), *c.conn, http.MethodPost, url, item)
}
