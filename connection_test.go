package solr

import (
	"net/http"
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

func TestNewRetryableConnection(t *testing.T) {
	_, err := NewRetryableConnection("", "mycore", http.DefaultClient, 500*time.Millisecond, false)
	if err == nil {
		t.Fatal("shouldn't run without a host")
	}

	_, err = NewRetryableConnection("invalid", "mycore", http.DefaultClient, 500*time.Millisecond, false)
	if err == nil {
		t.Fatal("shouldn't run without a proper host")
	}

	_, err = NewRetryableConnection("http://localhost", "", http.DefaultClient, 500*time.Millisecond, false)
	if err == nil {
		t.Fatal("shouldn't run without a core defined")
	}

	c, err := NewRetryableConnection("http://localhost:8983", "mycore", http.DefaultClient, 500*time.Millisecond, false)
	if err != nil {
		t.Fatal("shouldn't get an error but got one")
	}

	_, err = NewSingleClient(c)
	if err != nil {
		t.Fatal("shouldn't get an error but got one")
	}
}
