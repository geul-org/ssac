//ff:func feature=pkg-session type=util control=sequence
//ff:what 메모리 세션 Set — JSON 직렬화 후 맵에 저장
package session

import (
	"context"
	"encoding/json"
	"time"
)

func (s *memorySession) Set(_ context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.store[key] = memoryEntry{value: string(data), expiresAt: time.Now().Add(ttl)}
	s.mu.Unlock()
	return nil
}
