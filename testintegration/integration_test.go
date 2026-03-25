//go:build integration

package testintegration

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/mecenat/solr"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupSolr(t *testing.T) (solr.Client, func()) {
	t.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "solr:8.11",
		ExposedPorts: []string{"8983/tcp"},
		Cmd:          []string{"solr-precreate", "testcore"},
		WaitingFor:   wait.ForHTTP("/solr/testcore/admin/ping").WithPort("8983/tcp").WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start solr container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	port, err := container.MappedPort(ctx, "8983")
	if err != nil {
		t.Fatal(err)
	}

	solrURL := fmt.Sprintf("http://%s:%s", host, port.Port())

	conn, err := solr.NewConnection(solrURL, "testcore", solr.NewDefaultHTTPClient())
	if err != nil {
		t.Fatal(err)
	}
	client, err := solr.NewSingleClient(conn)
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		container.Terminate(ctx)
	}

	return client, cleanup
}

func TestIntegrationSearchShortURL(t *testing.T) {
	client, cleanup := setupSolr(t)
	defer cleanup()

	ctx := context.Background()

	err := client.Ping(ctx)
	if err != nil {
		t.Fatalf("ping failed: %v", err)
	}

	// Short query — should use GET
	q := solr.NewQuery(nil)
	q.AddQuery("", "*:*")
	res, err := client.Search(ctx, q)
	if err != nil {
		t.Fatalf("short URL search failed: %v", err)
	}
	if res.Header == nil {
		t.Fatal("expected response header")
	}
}

func TestIntegrationSearchLongURL(t *testing.T) {
	client, cleanup := setupSolr(t)
	defer cleanup()

	ctx := context.Background()

	// Build a query with enough filter queries to push the URL over 2048 chars
	q := solr.NewQuery(nil)
	q.AddQuery("", "*:*")
	for i := 0; i < 200; i++ {
		q.AddParam(fmt.Sprintf("param%d", i), strings.Repeat("x", 10))
	}

	// This URL will be >2048 chars, triggering POST with x-www-form-urlencoded
	res, err := client.Search(ctx, q)
	if err != nil {
		t.Fatalf("long URL search failed: %v", err)
	}
	if res.Header == nil {
		t.Fatal("expected response header")
	}
}

func TestIntegrationCreateAndSearch(t *testing.T) {
	client, cleanup := setupSolr(t)
	defer cleanup()

	ctx := context.Background()

	// Create a document
	doc := map[string]interface{}{
		"id": "1",
	}
	_, err := client.Create(ctx, doc, &solr.WriteOptions{Commit: true})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Search for it with a short query (GET)
	q := solr.NewQuery(nil)
	q.AddQuery("id", "1")
	res, err := client.Search(ctx, q)
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if res.Data == nil || res.Data.NumFound != 1 {
		t.Fatalf("expected 1 result, got %+v", res.Data)
	}

	// Now search with a long query (POST) that still matches the same doc
	qLong := solr.NewQuery(nil)
	qLong.AddQuery("id", "1")
	// Add many params to push URL over 2048
	for i := 0; i < 200; i++ {
		qLong.AddParam(fmt.Sprintf("param%d", i), strings.Repeat("x", 10))
	}

	res, err = client.Search(ctx, qLong)
	if err != nil {
		t.Fatalf("long URL search after create failed: %v", err)
	}
	if res.Data == nil || res.Data.NumFound != 1 {
		t.Fatalf("expected 1 result from long URL search, got %+v", res.Data)
	}
}

func TestIntegrationBatchGetLongURL(t *testing.T) {
	client, cleanup := setupSolr(t)
	defer cleanup()

	ctx := context.Background()

	// Create a document first
	doc := map[string]interface{}{
		"id": "get-test-1",
	}
	_, err := client.Create(ctx, doc, &solr.WriteOptions{Commit: true})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// BatchGet with enough IDs to push URL over 2048
	ids := make([]string, 100)
	ids[0] = "get-test-1"
	for i := 1; i < 100; i++ {
		ids[i] = fmt.Sprintf("nonexistent-%s", strings.Repeat("x", 20))
	}

	res, err := client.BatchGet(ctx, ids, "")
	if err != nil {
		t.Fatalf("batch get with long URL failed: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil response")
	}
}
