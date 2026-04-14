//ff:func feature=pkg-session type=util control=sequence
//ff:what PostgreSQL 세션 Delete — 키 기반 행 삭제
package session

import "context"

func (s *postgresSession) Delete(ctx context.Context, key string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM fullend_sessions WHERE key = $1`, key)
	return err
}
