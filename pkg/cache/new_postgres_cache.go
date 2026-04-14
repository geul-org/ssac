//ff:func feature=pkg-cache type=loader control=sequence
//ff:what PostgreSQL 캐시 생성 — fullend_cache 테이블 자동 생성 후 인스턴스 반환
package cache

import (
	"context"
	"database/sql"
)

// NewPostgresCache creates a CacheModel backed by PostgreSQL.
// It auto-creates the fullend_cache table if not exists.
func NewPostgresCache(ctx context.Context, db *sql.DB) (CacheModel, error) {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS fullend_cache (
			key        TEXT PRIMARY KEY,
			value      TEXT NOT NULL,
			expires_at TIMESTAMPTZ NOT NULL
		)`)
	if err != nil {
		return nil, err
	}
	return &postgresCache{db: db}, nil
}
