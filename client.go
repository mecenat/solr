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

func (c *Client) Search(q *Query) (*Response, error) {
	url := c.BasePath + "/select?" + q.String()
	return c.request(context.Background(), http.MethodGet, url, nil)
}

// Get performs a real time get that returns the latest version of the specified document(s)
// https://lucene.apache.org/solr/guide/8_5/realtime-get.html
func (c *Client) Get(id string) (*Response, error) {
	query := fmt.Sprintf("?id=%s", id)
	url := c.BasePath + "/get" + query
	return c.request(context.Background(), http.MethodGet, url, nil)
}

// BatchGet performs a real time get that returns the latest version of the specified document(s)
// https://lucene.apache.org/solr/guide/8_5/realtime-get.html
func (c *Client) BatchGet(ids []string, filter string) (*Response, error) {
	query := fmt.Sprintf("?ids=%s&fq=%s", strings.Join(ids, ","), url.QueryEscape(filter))
	url := c.BasePath + "/get" + query
	return c.request(context.Background(), http.MethodGet, url, nil)
}

func (c *Client) Create(item interface{}, opts *WriteOptions) (*Response, error) {
	url := c.formatURL("/update/json/docs", opts.formatQueryFromOpts())

	itemBytes, err := interfaceToBytes(item)
	if err != nil {
		return nil, err
	}

	err = isJSON(itemBytes)
	if err != nil {
		return nil, fmt.Errorf("Invalid JSON provided: %s", err)
	}

	return c.request(context.Background(), http.MethodPost, url, itemBytes)
}

func (c *Client) BatchCreate(items interface{}) (*Response, error) {
	url := c.BasePath + "/update/json"

	itemBytes, err := interfaceToBytes(items)
	if err != nil {
		return nil, err
	}

	err = isArrayOfJSON(itemBytes)
	if err != nil {
		return nil, fmt.Errorf("Invalid JSON provided: %s", err)
	}

	return c.request(context.Background(), http.MethodPost, url, itemBytes)
}

func (c *Client) Commit() (*Response, error) {
	url := c.BasePath + "/update?commit=true"
	return c.request(context.Background(), http.MethodPost, url, nil)
}

func (c *Client) DeleteByID(id string) (*Response, error) {
	path := c.BasePath + "/update"

	qq := map[string]map[string]interface{}{"delete": {"id": id}}

	bdBytes, err := interfaceToBytes(qq)
	if err != nil {
		return nil, err
	}

	return c.request(context.Background(), http.MethodPost, path, bdBytes)
}

func (c *Client) DeleteByQuery(query string) (*Response, error) {
	path := c.BasePath + "/update"

	qq := map[string]map[string]interface{}{"delete": {"query": url.QueryEscape(query)}}

	bdBytes, err := interfaceToBytes(qq)
	if err != nil {
		return nil, err
	}

	return c.request(context.Background(), http.MethodPost, path, bdBytes)
}
