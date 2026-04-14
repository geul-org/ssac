//ff:func feature=pkg-session type=test control=sequence
//ff:what MemorySession의 TTL 만료 동작을 검증한다
package session

import (
	"context"
	"testing"
	"time"
)

func TestMemorySession_Expiry(t *testing.T) {
	s := NewMemorySession()
	ctx := context.Background()

	if err := s.Set(ctx, "k2", "data", 1*time.Millisecond); err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Millisecond)

	val, err := s.Get(ctx, "k2")
	if err != nil {
		t.Fatal(err)
	}
	if val != "" {
		t.Errorf("expected empty after expiry, got %q", val)
	}
}
