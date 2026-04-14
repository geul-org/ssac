//ff:func feature=pkg-cache type=test control=sequence
//ff:what MemoryCache의 Set/Get/Delete 기본 동작을 검증한다
package cache

import (
	"context"
	"testing"
	"time"
)

func TestMemoryCache_SetGetDelete(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	if err := c.Set(ctx, "k1", "hello", 10*time.Second); err != nil {
		t.Fatal(err)
	}

	val, err := c.Get(ctx, "k1")
	if err != nil {
		t.Fatal(err)
	}
	if val != `"hello"` {
		t.Errorf("expected %q, got %q", `"hello"`, val)
	}

	if err := c.Delete(ctx, "k1"); err != nil {
		t.Fatal(err)
	}

	val, err = c.Get(ctx, "k1")
	if err != nil {
		t.Fatal(err)
	}
	if val != "" {
		t.Errorf("expected empty after delete, got %q", val)
	}
}
