//ff:func feature=pkg-auth type=test control=sequence topic=auth-jwt
//ff:what IssueToken → VerifyToken 라운드트립으로 claims가 보존되는지 검증한다
package auth

import (
	"encoding/json"
	"testing"
	"time"
)

// sampleClaim mirrors a project-local CurrentUser with JSON tags matching
// claim keys expected by the application. auth.IssueToken marshals this
// struct into jwt.MapClaims via encoding/json.
type sampleClaim struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	OrgID  int64  `json:"org_id"`
}

func TestIssueVerifyRoundTrip(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-test-secret-test-secret-xx")
	Configure(Config{SecretEnv: "JWT_SECRET", AccessTTL: time.Minute, RefreshTTL: time.Hour})

	src := sampleClaim{UserID: 1, Email: "a@b.c", Role: "admin", OrgID: 42}
	issued, err := IssueToken(IssueTokenRequest{Claims: src})
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}
	if issued.AccessToken == "" {
		t.Fatal("empty access token")
	}
	if issued.ExpiresAt.Before(time.Now()) {
		t.Fatalf("ExpiresAt in the past: %v", issued.ExpiresAt)
	}

	verified, err := VerifyToken(VerifyTokenRequest{Token: issued.AccessToken})
	if err != nil {
		t.Fatalf("VerifyToken: %v", err)
	}

	// Round-trip MapClaims back into sampleClaim and compare.
	raw, err := json.Marshal(verified.Claims)
	if err != nil {
		t.Fatalf("marshal claims: %v", err)
	}
	var got sampleClaim
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal claims: %v", err)
	}
	if got != src {
		t.Fatalf("claims mismatch\nwant: %+v\ngot:  %+v", src, got)
	}
}

func TestRefreshTokenRoundTrip(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-test-secret-test-secret-xx")
	Configure(Config{SecretEnv: "JWT_SECRET", AccessTTL: time.Minute, RefreshTTL: time.Hour})

	src := sampleClaim{UserID: 7, Email: "r@s.t", Role: "user", OrgID: 9}
	issued, err := RefreshToken(RefreshTokenRequest{Claims: src})
	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}
	if issued.RefreshToken == "" {
		t.Fatal("empty refresh token")
	}

	verified, err := VerifyToken(VerifyTokenRequest{Token: issued.RefreshToken})
	if err != nil {
		t.Fatalf("VerifyToken (refresh): %v", err)
	}
	if int64(verified.Claims["user_id"].(float64)) != 7 {
		t.Fatalf("user_id mismatch: %v", verified.Claims["user_id"])
	}
}

func TestIssueToken_SecretEnvEmpty(t *testing.T) {
	t.Setenv("JWT_SECRET", "")
	Configure(Config{SecretEnv: "JWT_SECRET", AccessTTL: time.Minute})

	if _, err := IssueToken(IssueTokenRequest{Claims: sampleClaim{UserID: 1}}); err == nil {
		t.Fatal("expected error when JWT_SECRET is empty")
	}
}

func TestIssueToken_SecretEnvNameNotConfigured(t *testing.T) {
	Configure(Config{SecretEnv: "", AccessTTL: time.Minute})
	if _, err := IssueToken(IssueTokenRequest{Claims: sampleClaim{UserID: 1}}); err == nil {
		t.Fatal("expected error when SecretEnv is not configured")
	}
}

func TestVerifyToken_RejectsNonHS256(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-test-secret-test-secret-xx")
	Configure(Config{SecretEnv: "JWT_SECRET", AccessTTL: time.Minute})

	// Header '{"alg":"none","typ":"JWT"}' base64url: eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0
	// Payload '{"user_id":1}' base64url: eyJ1c2VyX2lkIjoxfQ
	// Signature empty. Forged alg=none token must be rejected.
	tok := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJ1c2VyX2lkIjoxfQ."
	if _, err := VerifyToken(VerifyTokenRequest{Token: tok}); err == nil {
		t.Fatal("expected error when non-HS256 token is presented")
	}
}

func TestClaimsToMap_Nil(t *testing.T) {
	m, err := claimsToMap(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 0 {
		t.Fatalf("expected empty map, got %v", m)
	}
}
