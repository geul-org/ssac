//ff:func feature=pkg-cache type=test control=sequence
//ff:what MemoryCache의 TTL 만료 동작을 검증한다
package cache

import (
	"context"
	"testing"
	"time"
)

func TestMemoryCache_Expiry(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	if err := c.Set(ctx, "k2", "data", 1*time.Millisecond); err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Millisecond)

	val, err := c.Get(ctx, "k2")
	if err != nil {
		t.Fatal(err)
	}
	if val != "" {
		t.Errorf("expected empty after expiry, got %q", val)
	}
}
