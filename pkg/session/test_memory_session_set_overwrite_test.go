//ff:func feature=pkg-session type=test control=sequence
//ff:what MemorySession의 동일 키 덮어쓰기 동작을 검증한다
package session

import (
	"context"
	"testing"
	"time"
)

func TestMemorySession_SetOverwrite(t *testing.T) {
	s := NewMemorySession()
	ctx := context.Background()

	s.Set(ctx, "k3", "v1", 10*time.Second)
	s.Set(ctx, "k3", "v2", 10*time.Second)

	val, _ := s.Get(ctx, "k3")
	if val != `"v2"` {
		t.Errorf("expected %q, got %q", `"v2"`, val)
	}
}
