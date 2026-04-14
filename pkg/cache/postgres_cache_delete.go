//ff:func feature=pkg-cache type=util control=sequence
//ff:what PostgreSQL 캐시 Delete — 키 기반 행 삭제
package cache

import "context"

func (c *postgresCache) Delete(ctx context.Context, key string) error {
	_, err := c.db.ExecContext(ctx, `DELETE FROM fullend_cache WHERE key = $1`, key)
	return err
}
