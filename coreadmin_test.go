package solr

import (
	"context"
	"net/http"
	"testing"
)

func TestNewCoreAdminInvalidUrl(t *testing.T) {
	_, err := NewCoreAdmin(context.Background(), "", http.DefaultClient)
	if err == nil {
		t.Fatal("shouldn't run without a host")
	}

	_, err = NewCoreAdmin(context.Background(), "invalid", http.DefaultClient)
	if err == nil {
		t.Fatal("shouldn't run without a proper host")
	}
}

func TestSplitParams(t *testing.T) {
	ctx := context.Background()
	ca, err := NewCoreAdmin(ctx, "http://localhost:8983", http.DefaultClient)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	splitOpts := &CoreSplitOpts{
		Path:       []string{"path1", "path2"},
		TargetCore: []string{"core1", "core2"},
		Ranges:     "",
		SplitKey:   "",
		AsyncID:    "",
	}

	_, err = ca.Split(ctx, "test", splitOpts)
	if err == nil {
		t.Fatal("shouldn't be possible to run with both path & targetCore")
	}
}

func TestSplitBothRangesAndKey(t *testing.T) {
	ctx := context.Background()
	ca, err := NewCoreAdmin(ctx, "http://localhost:8983", http.DefaultClient)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	splitOpts := &CoreSplitOpts{
		Ranges:   "ranges",
		SplitKey: "key",
	}

	_, err = ca.Split(ctx, "test", splitOpts)
	if err == nil {
		t.Fatal("shouldn't be possible to run with both ranges & splitKey")
	}
}
