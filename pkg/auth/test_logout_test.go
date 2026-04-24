//ff:func feature=pkg-auth type=test control=sequence topic=auth-refresh
//ff:what Logout idempotent 검증 — 미존재/이미 revoked token 도 nil error (memoryRefreshStore 기반)
package auth

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestLogout_IdempotentRevoke(t *testing.T) {
	store := NewMemoryRefreshStore()
	token := "plaintext.refresh.jwt"
	claims, _ := json.Marshal(map[string]any{"user_id": int64(1)})
	if err := store.Create(context.Background(), token, claims, time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}

	// First Logout revokes.
	out, err := Logout(context.Background(), store, token)
	if err != nil {
		t.Fatalf("Logout: %v", err)
	}
	if !out.Success {
		t.Fatalf("expected Success=true, got %+v", out)
	}
	// Second Logout on the same (now revoked) token is a no-op.
	out, err = Logout(context.Background(), store, token)
	if err != nil {
		t.Fatalf("Logout idempotent: %v", err)
	}
	if !out.Success {
		t.Fatalf("expected Success=true on idempotent logout, got %+v", out)
	}
}

func TestLogout_EmptyTokenSilent(t *testing.T) {
	store := NewMemoryRefreshStore()
	out, err := Logout(context.Background(), store, "")
	if err != nil {
		t.Fatalf("Logout empty: %v", err)
	}
	if !out.Success {
		t.Fatalf("empty token should still return Success=true, got %+v", out)
	}
}

func TestLogout_NilStore(t *testing.T) {
	_, err := Logout(context.Background(), nil, "some.token")
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}
