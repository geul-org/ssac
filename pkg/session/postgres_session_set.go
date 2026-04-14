//ff:func feature=pkg-session type=util control=sequence
//ff:what PostgreSQL 세션 Set — JSON 직렬화 후 UPSERT
package session

import (
	"context"
	"encoding/json"
	"time"
)

func (s *postgresSession) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	expiresAt := time.Now().Add(ttl)
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO fullend_sessions (key, value, expires_at) VALUES ($1, $2, $3)
		ON CONFLICT (key) DO UPDATE SET value = $2, expires_at = $3`,
		key, string(data), expiresAt)
	return err
}
