package solr

import (
	"context"
	"net/http"
	"testing"
)

func TestNewSingleClientInvalidUrl(t *testing.T) {
	_, err := NewSingleClient(context.Background(), "", "mycore", http.DefaultClient)
	if err == nil {
		t.Fatal("shouldn't run without a host")
	}

	_, err = NewSingleClient(context.Background(), "invalid", "mycore", http.DefaultClient)
	if err == nil {
		t.Fatal("shouldn't run without a proper host")
	}

	_, err = NewSingleClient(context.Background(), "http://localhost", "", http.DefaultClient)
	if err == nil {
		t.Fatal("shouldn't run without a core defined")
	}
}
