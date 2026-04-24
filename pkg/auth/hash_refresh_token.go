//ff:func feature=pkg-auth type=util control=sequence topic=auth-refresh
//ff:what hashRefreshToken — sha256 해시 계산 (RefreshStore 구현체 공용)
package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashRefreshToken returns the sha256 hex digest of a refresh-token string.
// Shared by RefreshStore implementations so the DB never stores plaintext.
// Exported so yongol-generated postgres RefreshStore implementations can
// hash tokens identically to the memory reference implementation.
func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// MarshalClaimsJSON normalizes passthrough claims to a JSON blob. nil claims
// becomes "{}" so JSONB columns never receive NULL. Exported for reuse by
// yongol-generated RefreshStore implementations.
func MarshalClaimsJSON(in any) ([]byte, error) {
	if in == nil {
		return []byte(`{}`), nil
	}
	if raw, ok := in.([]byte); ok {
		return raw, nil
	}
	// json.RawMessage is []byte under the hood; the type assertion above
	// covers the common case. For other types the standard marshal path
	// applies.
	return marshalJSON(in)
}
