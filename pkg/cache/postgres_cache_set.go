//ff:func feature=pkg-cache type=util control=sequence
//ff:what PostgreSQL 캐시 Set — JSON 직렬화 후 UPSERT
package cache

import (
	"context"
	"encoding/json"
	"time"
)

func (c *postgresCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	expiresAt := time.Now().Add(ttl)
	_, err = c.db.ExecContext(ctx, `
		INSERT INTO fullend_cache (key, value, expires_at) VALUES ($1, $2, $3)
		ON CONFLICT (key) DO UPDATE SET value = $2, expires_at = $3`,
		key, string(data), expiresAt)
	return err
}
