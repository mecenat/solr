package solr

import (
	"context"
	"net/http"
	"testing"
)

func TestNewPrimaryReplicaClientInvalidUrl(t *testing.T) {
	_, err := NewPrimaryReplicaClient(context.Background(), "", "http://localhost", "mycore", http.DefaultClient, http.DefaultClient)
	if err == nil {
		t.Fatal("shouldn't run without a primary host")
	}

	_, err = NewPrimaryReplicaClient(context.Background(), "http://localhost", "", "mycore", http.DefaultClient, http.DefaultClient)
	if err == nil {
		t.Fatal("shouldn't run without a replica host")
	}

	_, err = NewPrimaryReplicaClient(context.Background(), "invalid", "http://localhost", "mycore", http.DefaultClient, http.DefaultClient)
	if err == nil {
		t.Fatal("shouldn't run without a proper primary host")
	}

	_, err = NewPrimaryReplicaClient(context.Background(), "http://localhost", "invalid", "mycore", http.DefaultClient, http.DefaultClient)
	if err == nil {
		t.Fatal("shouldn't run without a proper replica host")
	}

	_, err = NewPrimaryReplicaClient(context.Background(), "http://localhost", "http://localhos", "", http.DefaultClient, http.DefaultClient)
	if err == nil {
		t.Fatal("shouldn't run without a core defined")
	}
}
