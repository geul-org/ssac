//ff:type feature=pkg-auth type=store topic=auth-refresh
//ff:what memoryRefreshStore — 테스트/개발용 in-memory RefreshStore
package auth

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

// memoryRefreshStore keeps refresh-token rows in a map keyed by sha256 hash.
// Intended for tests and zero-config dev. Production deployments use the
// yongol-generated postgres RefreshStore backed by user sqlc Queries.
type memoryRefreshStore struct {
	mu   sync.Mutex
	rows map[string]*memoryRefreshRow
}

type memoryRefreshRow struct {
	claims    json.RawMessage
	expiresAt time.Time
	revokedAt *time.Time
}

// NewMemoryRefreshStore returns a RefreshStore backed by an in-memory map.
func NewMemoryRefreshStore() RefreshStore {
	return &memoryRefreshStore{rows: make(map[string]*memoryRefreshRow)}
}

func (s *memoryRefreshStore) Create(_ context.Context, token string, claims any, expiresAt time.Time) error {
	raw, err := MarshalClaimsJSON(claims)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rows[HashRefreshToken(token)] = &memoryRefreshRow{
		claims:    raw,
		expiresAt: expiresAt,
	}
	return nil
}

func (s *memoryRefreshStore) Consume(_ context.Context, token string) (json.RawMessage, error) {
	hash := HashRefreshToken(token)
	s.mu.Lock()
	defer s.mu.Unlock()
	row, ok := s.rows[hash]
	if !ok {
		return nil, ErrRefreshTokenNotFound
	}
	if row.revokedAt != nil {
		return row.claims, ErrRefreshTokenReused
	}
	if time.Now().After(row.expiresAt) {
		return nil, ErrRefreshTokenNotFound
	}
	now := time.Now()
	row.revokedAt = &now
	return row.claims, nil
}

func (s *memoryRefreshStore) Revoke(_ context.Context, token string) error {
	hash := HashRefreshToken(token)
	s.mu.Lock()
	defer s.mu.Unlock()
	if row, ok := s.rows[hash]; ok && row.revokedAt == nil {
		now := time.Now()
		row.revokedAt = &now
	}
	return nil
}

func (s *memoryRefreshStore) RevokeAll(_ context.Context, matcher ClaimMatcher) error {
	if len(matcher) == 0 {
		return ErrEmptyMatcher
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, row := range s.rows {
		if row.revokedAt != nil {
			continue
		}
		if matchClaims(row.claims, matcher) {
			now := time.Now()
			row.revokedAt = &now
		}
	}
	return nil
}

// matchClaims returns true when every key/value in matcher is present in
// the stored JSON claims (shallow equality on top-level keys).
func matchClaims(stored json.RawMessage, matcher ClaimMatcher) bool {
	var m map[string]any
	if err := json.Unmarshal(stored, &m); err != nil {
		return false
	}
	for k, want := range matcher {
		got, ok := m[k]
		if !ok {
			return false
		}
		// Compare via JSON serialization so numeric types coerce uniformly.
		a, _ := json.Marshal(got)
		b, _ := json.Marshal(want)
		if string(a) != string(b) {
			return false
		}
	}
	return true
}
