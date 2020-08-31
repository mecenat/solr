package solr

import (
	"context"
	"net/http"
	"testing"
)

func TestNewSchemaAPIInvalidUrl(t *testing.T) {
	_, err := NewSchemaAPI(context.Background(), "", "mycore", http.DefaultClient)
	if err == nil {
		t.Fatal("shouldn't run without a host")
	}

	_, err = NewSchemaAPI(context.Background(), "invalid", "mycore", http.DefaultClient)
	if err == nil {
		t.Fatal("shouldn't run without a proper host")
	}

	_, err = NewSchemaAPI(context.Background(), "http://localhost", "", http.DefaultClient)
	if err == nil {
		t.Fatal("shouldn't run without a core defined")
	}
}
