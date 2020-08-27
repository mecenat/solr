package solr

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
	Username   string
	Password   string
}

func (c *Connection) request(ctx context.Context, method, url string, body []byte) (*Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	res, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	var r Response
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	if r.Error != nil {
		return &r, r.Error
	}

	return &r, nil
}

func formatBasePath(host, core string) string {
	if strings.HasSuffix(host, "/solr") {
		return fmt.Sprintf("%s/%s", host, core)
	}
	return fmt.Sprintf("%s/solr/%s", host, core)
}

func formatDocEntry(doc Doc) map[string]interface{} {
	return map[string]interface{}{"doc": doc}
}

func formatDeleteByID(id string) Doc {
	return Doc{"id": id}
}

func formatDeleteByQuery(query string) Doc {
	return Doc{"query": query}
}

func isJSON(input []byte) error {
	var js map[string]interface{}
	return json.Unmarshal(input, &js)
}

func isArrayOfJSON(input []byte) error {
	var js []*map[string]interface{}
	return json.Unmarshal(input, &js)
}

func interfaceToBytes(a interface{}) ([]byte, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return b, err
}
