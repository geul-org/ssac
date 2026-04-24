//ff:func feature=pkg-auth type=test control=sequence topic=auth-refresh
//ff:what RefreshRotate happy path + reuse 감지 + invalid JWT — memoryRefreshStore 기반
package auth

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestRefreshRotate_HappyPath(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-test-secret-test-secret-xx")
	Configure(Config{SecretEnv: "JWT_SECRET", AccessTTL: time.Minute, RefreshTTL: time.Hour})

	src := sampleClaim{UserID: 1, Email: "a@b.c", Role: "admin", OrgID: 42}
	issued, err := RefreshToken(RefreshTokenRequest{Claims: src})
	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}
	rawClaims, _ := json.Marshal(src)

	store := NewMemoryRefreshStore()
	if err := store.Create(context.Background(), issued.RefreshToken, rawClaims, issued.ExpiresAt); err != nil {
		t.Fatalf("Create: %v", err)
	}

	out, err := RefreshRotate(context.Background(), store, issued.RefreshToken, false)
	if err != nil {
		t.Fatalf("RefreshRotate: %v", err)
	}
	if out.AccessToken == "" || out.RefreshToken == "" {
		t.Fatalf("missing tokens in response: %+v", out)
	}
	verified, err := VerifyToken(VerifyTokenRequest{Token: out.AccessToken})
	if err != nil {
		t.Fatalf("verify new access: %v", err)
	}
	if int64(verified.Claims["user_id"].(float64)) != 1 {
		t.Fatalf("claims not preserved in rotation: %v", verified.Claims)
	}
}

func TestRefreshRotate_ReuseDetected(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-test-secret-test-secret-xx")
	Configure(Config{SecretEnv: "JWT_SECRET", AccessTTL: time.Minute, RefreshTTL: time.Hour})

	src := sampleClaim{UserID: 1, Email: "a@b.c", Role: "admin", OrgID: 42}
	issued, err := RefreshToken(RefreshTokenRequest{Claims: src})
	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}
	rawClaims, _ := json.Marshal(src)

	store := NewMemoryRefreshStore()
	ctx := context.Background()
	if err := store.Create(ctx, issued.RefreshToken, rawClaims, issued.ExpiresAt); err != nil {
		t.Fatalf("Create: %v", err)
	}
	// Manually revoke the row so the next Consume surfaces reuse — matches
	// the production path where a stolen/leaked refresh token is presented
	// after the legitimate holder already rotated.
	if err := store.Revoke(ctx, issued.RefreshToken); err != nil {
		t.Fatalf("Revoke: %v", err)
	}
	// Attempting to rotate on the revoked token must surface reuse.
	_, err = RefreshRotate(ctx, store, issued.RefreshToken, true)
	if !errors.Is(err, ErrRefreshTokenReused) {
		t.Fatalf("expected ErrRefreshTokenReused, got %v", err)
	}
}

func TestRefreshRotate_InvalidJWT(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-test-secret-test-secret-xx")
	Configure(Config{SecretEnv: "JWT_SECRET", AccessTTL: time.Minute, RefreshTTL: time.Hour})

	store := NewMemoryRefreshStore()
	_, err := RefreshRotate(context.Background(), store, "garbage.not.jwt", false)
	if err == nil {
		t.Fatal("expected error for invalid JWT, got nil")
	}
	if !containsVerifyPrefix(err.Error()) {
		t.Fatalf("expected verify-prefixed error, got %q", err.Error())
	}
}

// containsVerifyPrefix is a local mirror of the previous refresh_handler
// helper so we can keep the error-wrap assertion without reintroducing the
// handler.
func containsVerifyPrefix(s string) bool {
	const prefix = "auth: verify refresh token:"
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}
