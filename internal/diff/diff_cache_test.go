package diff

import (
	"testing"
)

func makeCacheReport(addr string) *Report {
	return &Report{
		Entries: []ResourceChange{
			{Address: addr, ResourceType: "aws_s3_bucket", Action: ActionAdded},
		},
	}
}

func TestDiffCache_SetAndGet(t *testing.T) {
	c := NewDiffCache()
	r := makeCacheReport("aws_s3_bucket.example")
	c.Set("key1", r)

	got := c.Get("key1")
	if got == nil {
		t.Fatal("expected cached report, got nil")
	}
	if len(got.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got.Entries))
	}
}

func TestDiffCache_GetMiss(t *testing.T) {
	c := NewDiffCache()
	if got := c.Get("missing"); got != nil {
		t.Fatalf("expected nil for cache miss, got %+v", got)
	}
}

func TestDiffCache_SetNilNoOp(t *testing.T) {
	c := NewDiffCache()
	c.Set("key1", nil)
	if c.Len() != 0 {
		t.Fatalf("expected 0 entries after nil set, got %d", c.Len())
	}
}

func TestDiffCache_Len(t *testing.T) {
	c := NewDiffCache()
	c.Set("a", makeCacheReport("addr.a"))
	c.Set("b", makeCacheReport("addr.b"))
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}

func TestDiffCache_Flush(t *testing.T) {
	c := NewDiffCache()
	c.Set("a", makeCacheReport("addr.a"))
	c.Flush()
	if c.Len() != 0 {
		t.Fatalf("expected 0 after flush, got %d", c.Len())
	}
}

func TestCacheKey_Deterministic(t *testing.T) {
	type payload struct{ Name string }
	k1, err := CacheKey(payload{"foo"}, payload{"bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	k2, err := CacheKey(payload{"foo"}, payload{"bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if k1 != k2 {
		t.Fatalf("expected identical keys, got %q and %q", k1, k2)
	}
}

func TestCacheKey_DifferentInputs(t *testing.T) {
	type payload struct{ Name string }
	k1, _ := CacheKey(payload{"foo"}, payload{"bar"})
	k2, _ := CacheKey(payload{"foo"}, payload{"baz"})
	if k1 == k2 {
		t.Fatal("expected different keys for different inputs")
	}
}
