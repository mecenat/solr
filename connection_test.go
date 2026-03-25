package solr

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewConnection(t *testing.T) {
	_, err := NewConnection("", "mycore", http.DefaultClient)
	if err == nil {
		t.Fatal("shouldn't run without a host")
	}

	_, err = NewConnection("invalid", "mycore", http.DefaultClient)
	if err == nil {
		t.Fatal("shouldn't run without a proper host")
	}

	_, err = NewConnection("http://localhost", "", http.DefaultClient)
	if err == nil {
		t.Fatal("shouldn't run without a core defined")
	}

	c, err := NewConnection("http://localhost:8983", "mycore", http.DefaultClient)
	if err != nil {
		t.Fatal("shouldn't get an error but got one")
	}

	_, err = NewSingleClient(c)
	if err != nil {
		t.Fatal("shouldn't get an error but got one")
	}
}

func TestConnectionDrainsResponseBody(t *testing.T) {
	var connCount int32
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"responseHeader":{"status":0,"QTime":1}}` + "\n\n\n"))
	}))
	ts.Config.ConnState = func(conn net.Conn, state http.ConnState) {
		if state == http.StateNew {
			atomic.AddInt32(&connCount, 1)
		}
	}
	ts.Start()
	defer ts.Close()

	client := &http.Client{Transport: &http.Transport{}}
	conn, err := NewConnection(ts.URL, "testcore", client)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	url := conn.formatBasePath() + "/select?q=*%3A*&wt=json"

	for i := 0; i < 3; i++ {
		_, err = conn.request(ctx, http.MethodGet, url, "application/json", nil)
		if err != nil {
			t.Fatal(err)
		}
	}

	if c := atomic.LoadInt32(&connCount); c != 1 {
		t.Errorf("expected 1 connection (reuse), got %d", c)
	}
}

func TestNewRetryableConnection(t *testing.T) {
	conf := &RetryableConfig{
		Timeout:      10 * time.Second,
		RetryWaitMin: 50 * time.Millisecond,
		RetryWaitMax: 2 * time.Second,
		RetryMax:     4,
		NoLog:        true,
	}

	_, err := NewRetryableConnection("", "mycore", http.DefaultClient, conf)
	if err == nil {
		t.Fatal("shouldn't run without a host")
	}

	_, err = NewRetryableConnection("invalid", "mycore", http.DefaultClient, conf)
	if err == nil {
		t.Fatal("shouldn't run without a proper host")
	}

	_, err = NewRetryableConnection("http://localhost", "", http.DefaultClient, conf)
	if err == nil {
		t.Fatal("shouldn't run without a core defined")
	}

	c, err := NewRetryableConnection("http://localhost:8983", "mycore", http.DefaultClient, conf)
	if err != nil {
		t.Fatal("shouldn't get an error but got one")
	}

	_, err = NewSingleClient(c)
	if err != nil {
		t.Fatal("shouldn't get an error but got one")
	}
}

func TestNewDefaultHTTPClient(t *testing.T) {
	client := NewDefaultHTTPClient()
	if client == nil {
		t.Fatal("expected non-nil client")
	}

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("expected *http.Transport")
	}

	if transport.MaxIdleConns != 100 {
		t.Fatalf("expected MaxIdleConns=100, got %d", transport.MaxIdleConns)
	}
	if transport.MaxIdleConnsPerHost != 100 {
		t.Fatalf("expected MaxIdleConnsPerHost=100, got %d", transport.MaxIdleConnsPerHost)
	}
	if transport.IdleConnTimeout != 90*time.Second {
		t.Fatalf("expected IdleConnTimeout=90s, got %v", transport.IdleConnTimeout)
	}
	if !transport.ForceAttemptHTTP2 {
		t.Fatal("expected ForceAttemptHTTP2=true")
	}
}

func TestNewDefaultHTTPClientWorksWithNewConnection(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"responseHeader":{"status":0,"QTime":1}}`))
	}))
	defer ts.Close()

	client := NewDefaultHTTPClient()
	conn, err := NewConnection(ts.URL, "testcore", client)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	url := conn.formatBasePath() + "/admin/ping"
	_, err = conn.request(ctx, http.MethodGet, url, "application/json", nil)
	if err != nil {
		t.Fatalf("request with default client failed: %v", err)
	}
}
