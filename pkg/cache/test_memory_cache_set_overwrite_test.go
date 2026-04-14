//ff:func feature=pkg-cache type=test control=sequence
//ff:what MemoryCache의 동일 키 덮어쓰기 동작을 검증한다
package cache

import (
	"context"
	"testing"
	"time"
)

func TestMemoryCache_SetOverwrite(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	c.Set(ctx, "k3", "v1", 10*time.Second)
	c.Set(ctx, "k3", "v2", 10*time.Second)

	val, _ := c.Get(ctx, "k3")
	if val != `"v2"` {
		t.Errorf("expected %q, got %q", `"v2"`, val)
	}
}
