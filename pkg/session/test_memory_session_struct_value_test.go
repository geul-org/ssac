//ff:func feature=pkg-session type=test control=sequence
//ff:what MemorySession에 구조체 값을 저장하고 조회하는 동작을 검증한다
package session

import (
	"context"
	"testing"
	"time"
)

func TestMemorySession_StructValue(t *testing.T) {
	s := NewMemorySession()
	ctx := context.Background()

	data := map[string]string{"user": "alice", "role": "admin"}
	if err := s.Set(ctx, "k4", data, 10*time.Second); err != nil {
		t.Fatal(err)
	}

	val, err := s.Get(ctx, "k4")
	if err != nil {
		t.Fatal(err)
	}
	if val == "" {
		t.Error("expected non-empty value for struct")
	}
}
