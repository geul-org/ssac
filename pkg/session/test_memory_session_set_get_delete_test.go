//ff:func feature=pkg-session type=test control=sequence
//ff:what MemorySession의 Set/Get/Delete 기본 동작을 검증한다
package session

import (
	"context"
	"testing"
	"time"
)

func TestMemorySession_SetGetDelete(t *testing.T) {
	s := NewMemorySession()
	ctx := context.Background()

	if err := s.Set(ctx, "k1", "hello", 10*time.Second); err != nil {
		t.Fatal(err)
	}

	val, err := s.Get(ctx, "k1")
	if err != nil {
		t.Fatal(err)
	}
	if val != `"hello"` {
		t.Errorf("expected %q, got %q", `"hello"`, val)
	}

	if err := s.Delete(ctx, "k1"); err != nil {
		t.Fatal(err)
	}

	val, err = s.Get(ctx, "k1")
	if err != nil {
		t.Fatal(err)
	}
	if val != "" {
		t.Errorf("expected empty after delete, got %q", val)
	}
}
