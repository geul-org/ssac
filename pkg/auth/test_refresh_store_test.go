//ff:func feature=pkg-auth type=test control=sequence topic=auth-refresh
//ff:what memoryRefreshStore Create/Consume/RevokeAll/재사용 감지 동작 검증
package auth

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestMemoryRefreshStore_CreateConsume(t *testing.T) {
	store := NewMemoryRefreshStore()
	ctx := context.Background()
	token := "plaintext.refresh.jwt"
	claims := sampleClaim{UserID: 1, Email: "a@b.c", Role: "admin", OrgID: 42}
	expiresAt := time.Now().Add(time.Hour)

	if err := store.Create(ctx, token, claims, expiresAt); err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := store.Consume(ctx, token)
	if err != nil {
		t.Fatalf("Consume: %v", err)
	}
	var back sampleClaim
	if err := json.Unmarshal(got, &back); err != nil {
		t.Fatalf("unmarshal claims: %v", err)
	}
	if back != claims {
		t.Fatalf("claims mismatch\nwant: %+v\ngot:  %+v", claims, back)
	}
}

func TestMemoryRefreshStore_ConsumeReuseDetected(t *testing.T) {
	store := NewMemoryRefreshStore()
	ctx := context.Background()
	token := "reused.refresh.jwt"
	claims := sampleClaim{UserID: 1}
	if err := store.Create(ctx, token, claims, time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	// First consume succeeds.
	if _, err := store.Consume(ctx, token); err != nil {
		t.Fatal(err)
	}
	// Second consume must surface reuse with the revoked row's claims.
	claimsBack, err := store.Consume(ctx, token)
	if !errors.Is(err, ErrRefreshTokenReused) {
		t.Fatalf("expected ErrRefreshTokenReused, got %v", err)
	}
	if len(claimsBack) == 0 {
		t.Fatal("expected revoked-row claims to be returned on reuse")
	}
}

func TestMemoryRefreshStore_ConsumeMissing(t *testing.T) {
	store := NewMemoryRefreshStore()
	if _, err := store.Consume(context.Background(), "missing.jwt"); !errors.Is(err, ErrRefreshTokenNotFound) {
		t.Fatalf("expected ErrRefreshTokenNotFound, got %v", err)
	}
}

func TestMemoryRefreshStore_RevokeAllRejectsEmptyMatcher(t *testing.T) {
	store := NewMemoryRefreshStore()
	if err := store.RevokeAll(context.Background(), ClaimMatcher{}); err == nil {
		t.Fatal("expected empty matcher to be rejected")
	}
}

func TestMemoryRefreshStore_RevokeAllWithMatcher(t *testing.T) {
	store := NewMemoryRefreshStore()
	ctx := context.Background()
	// Seed three tokens; two share user_id=1, one has user_id=2.
	seed := []struct {
		tok    string
		claims sampleClaim
	}{
		{"a", sampleClaim{UserID: 1}},
		{"b", sampleClaim{UserID: 1}},
		{"c", sampleClaim{UserID: 2}},
	}
	for _, s := range seed {
		if err := store.Create(ctx, s.tok, s.claims, time.Now().Add(time.Hour)); err != nil {
			t.Fatal(err)
		}
	}
	if err := store.RevokeAll(ctx, ClaimMatcher{"user_id": int64(1)}); err != nil {
		t.Fatalf("RevokeAll: %v", err)
	}
	// Tokens a, b should now be revoked.
	for _, tok := range []string{"a", "b"} {
		if _, err := store.Consume(ctx, tok); !errors.Is(err, ErrRefreshTokenReused) {
			t.Errorf("token %q: expected ErrRefreshTokenReused after RevokeAll, got %v", tok, err)
		}
	}
	// Token c should still consume cleanly.
	if _, err := store.Consume(ctx, "c"); err != nil {
		t.Errorf("token c: expected clean consume, got %v", err)
	}
}
