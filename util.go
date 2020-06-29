package solr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

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

func request(ctx context.Context, client *http.Client, method, url string, body []byte) (*Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req.WithContext(ctx))
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
		return nil, r.Error
	}

	return &r, nil
}
