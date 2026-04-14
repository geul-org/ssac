//ff:func feature=pkg-session type=loader control=sequence
//ff:what PostgreSQL 세션 생성 — fullend_sessions 테이블 자동 생성 후 인스턴스 반환
package session

import (
	"context"
	"database/sql"
)

// NewPostgresSession creates a SessionModel backed by PostgreSQL.
// It auto-creates the fullend_sessions table if not exists.
func NewPostgresSession(ctx context.Context, db *sql.DB) (SessionModel, error) {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS fullend_sessions (
			key        TEXT PRIMARY KEY,
			value      TEXT NOT NULL,
			expires_at TIMESTAMPTZ NOT NULL
		)`)
	if err != nil {
		return nil, err
	}
	return &postgresSession{db: db}, nil
}
