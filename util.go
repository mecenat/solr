package solr

import (
	"bytes"
	"context"
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

func request(ctx context.Context, conn Connection, method, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	return conn.httpClient.Do(req.WithContext(ctx))
}
