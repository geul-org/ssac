//ff:func feature=pkg-cache type=util control=selection
//ff:what PostgreSQL 캐시 Get — 만료되지 않은 값 조회, ErrNoRows 시 빈 문자열 반환
package cache

import (
	"context"
	"database/sql"
)

func (c *postgresCache) Get(ctx context.Context, key string) (string, error) {
	var value string
	err := c.db.QueryRowContext(ctx, `
		SELECT value FROM fullend_cache WHERE key = $1 AND expires_at > NOW()`,
		key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}
