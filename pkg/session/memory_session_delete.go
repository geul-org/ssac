//ff:func feature=pkg-session type=util control=sequence
//ff:what 메모리 세션 Delete — 맵에서 키 제거
package session

import "context"

func (s *memorySession) Delete(_ context.Context, key string) error {
	s.mu.Lock()
	delete(s.store, key)
	s.mu.Unlock()
	return nil
}
