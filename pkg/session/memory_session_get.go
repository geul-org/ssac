//ff:func feature=pkg-session type=util control=selection
//ff:what 메모리 세션 Get — 키 존재 여부와 만료 시각 확인 후 값 반환
package session

import (
	"context"
	"time"
)

func (s *memorySession) Get(_ context.Context, key string) (string, error) {
	s.mu.RLock()
	entry, ok := s.store[key]
	s.mu.RUnlock()
	if !ok || time.Now().After(entry.expiresAt) {
		return "", nil
	}
	return entry.value, nil
}
