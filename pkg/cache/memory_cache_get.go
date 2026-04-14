//ff:func feature=pkg-cache type=util control=selection
//ff:what 메모리 캐시 Get — 키 존재 여부와 만료 시각 확인 후 값 반환
package cache

import (
	"context"
	"time"
)

func (c *memoryCache) Get(_ context.Context, key string) (string, error) {
	c.mu.RLock()
	entry, ok := c.store[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(entry.expiresAt) {
		return "", nil
	}
	return entry.value, nil
}
