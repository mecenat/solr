package solr

import (
	"context"
	"encoding/json"
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

func (c *Client) Ping() (int, error) {
	url := c.BasePath + "/admin/ping"
	res, err := request(context.Background(), *c.conn, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	return res.StatusCode, nil
}

type GetResponse struct {
	Doc   *Doc           `json:"doc,omitempty"`
	Data  *ResponseData  `json:"response,omitempty"`
	Error *ResponseError `json:"error"`
}

type SearchResponse struct {
	Header *ResponseHeader         `json:"responseHeader"`
	Data   *ResponseData           `json:"response"`
	Error  *ResponseError          `json:"error"`
	Debug  *map[string]interface{} `json:"debug,omitempty"`
}

type ResponseError struct {
	Code    int64    `json:"code"`
	Message string   `json:"msg"`
	Meta    []string `json:"metadata"`
}

func (r *ResponseError) Error() string {
	return r.Message
}

type ResponseHeader struct {
	Status int64  `json:"status"`
	QTime  int64  `json:"QTime"`
	Params *Query `json:"params"`
}

type ResponseData struct {
	NumFound int64 `json:"numFound"`
	Start    int64 `json:"start"`
	Docs     Docs  `json:"docs"`
}

type Docs []*Doc
type Doc map[string]interface{}

func (d Docs) ToBytes() ([]byte, error) {
	return interfaceToBytes(d)
}

func (d *Doc) ToBytes() ([]byte, error) {
	return interfaceToBytes(d)
}

func interfaceToBytes(a interface{}) ([]byte, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return b, err
}

func (c *Client) Search(q *Query) (*SearchResponse, error) {
	url := c.BasePath + "/select?" + q.String()
	res, err := request(context.Background(), *c.conn, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	var sr SearchResponse
	defer res.Body.Close()

	// bd, err := ioutil.ReadAll(res.Body)

	err = json.NewDecoder(res.Body).Decode(&sr)
	if err != nil {
		return nil, err
	}

	if sr.Error != nil {
		return nil, fmt.Errorf("%s", sr.Error.Message)
	}

	// fmt.Println(string(bd))

	return &sr, nil
}

// Get performs a real time get that returns the latest version of the specified document(s)
// https://lucene.apache.org/solr/guide/8_5/realtime-get.html
func (c *Client) Get(ids []string, filter string) (*GetResponse, error) {
	query := fmt.Sprintf("?ids=%s&fq=%s", strings.Join(ids, ","), url.QueryEscape(filter))
	url := c.BasePath + "/get" + query
	res, err := request(context.Background(), *c.conn, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var gr GetResponse
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&gr)
	if err != nil {
		return nil, err
	}

	if gr.Error != nil {
		return nil, gr.Error
	}

	return &gr, nil
}
