package solr

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadUsesGetForShortURLs(t *testing.T) {
	var method string
	var contentType string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		contentType = r.Header.Get("Content-Type")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"responseHeader":{"status":0,"QTime":1}}`))
	}))
	defer ts.Close()

	conn, err := NewConnection(ts.URL, "testcore", &http.Client{})
	if err != nil {
		t.Fatal(err)
	}

	url := conn.formatBasePath() + "/select?q=*%3A*&wt=json"
	_, err = read(context.Background(), conn, url)
	if err != nil {
		t.Fatal(err)
	}

	if method != http.MethodGet {
		t.Errorf("expected GET for short URL, got %s", method)
	}
	if contentType != "application/json" {
		t.Errorf("expected application/json content type, got %s", contentType)
	}
}

func TestReadUsesPostForLongURLs(t *testing.T) {
	var method string
	var contentType string
	var body string
	var requestPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		contentType = r.Header.Get("Content-Type")
		requestPath = r.URL.Path
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"responseHeader":{"status":0,"QTime":1}}`))
	}))
	defer ts.Close()

	conn, err := NewConnection(ts.URL, "testcore", &http.Client{})
	if err != nil {
		t.Fatal(err)
	}

	// Build a URL longer than 2048 characters
	longValue := strings.Repeat("x", 2100)
	url := conn.formatBasePath() + "/select?q=" + longValue + "&wt=json"

	_, err = read(context.Background(), conn, url)
	if err != nil {
		t.Fatal(err)
	}

	if method != http.MethodPost {
		t.Errorf("expected POST for long URL, got %s", method)
	}
	if contentType != "application/x-www-form-urlencoded" {
		t.Errorf("expected application/x-www-form-urlencoded, got %s", contentType)
	}
	if !strings.HasSuffix(requestPath, "/select") {
		t.Errorf("expected request path to end with /select, got %s", requestPath)
	}
	if !strings.Contains(body, "q="+longValue) {
		truncated := body
		if len(truncated) > 100 {
			truncated = truncated[:100]
		}
		t.Errorf("expected body to contain query params, got %s", truncated)
	}
}
