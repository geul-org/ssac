//ff:func feature=pkg-cache type=util control=sequence
//ff:what 메모리 캐시 Delete — 맵에서 키 제거
package cache

import "context"

func (c *memoryCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	delete(c.store, key)
	c.mu.Unlock()
	return nil
}
