//ff:func feature=pkg-cache type=util control=sequence
//ff:what 메모리 캐시 Set — JSON 직렬화 후 맵에 저장
package cache

import (
	"context"
	"encoding/json"
	"time"
)

func (c *memoryCache) Set(_ context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	c.mu.Lock()
	c.store[key] = memoryEntry{value: string(data), expiresAt: time.Now().Add(ttl)}
	c.mu.Unlock()
	return nil
}
