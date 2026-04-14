//ff:func feature=pkg-session type=util control=selection
//ff:what PostgreSQL 세션 Get — 만료되지 않은 값 조회, ErrNoRows 시 빈 문자열 반환
package session

import (
	"context"
	"database/sql"
)

func (s *postgresSession) Get(ctx context.Context, key string) (string, error) {
	var value string
	err := s.db.QueryRowContext(ctx, `
		SELECT value FROM fullend_sessions WHERE key = $1 AND expires_at > NOW()`,
		key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}
