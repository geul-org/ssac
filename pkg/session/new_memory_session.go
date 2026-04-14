//ff:func feature=pkg-session type=loader control=sequence
//ff:what 인메모리 세션 생성 — 재시작 시 데이터 소멸
package session

// NewMemorySession creates an in-memory SessionModel. Data is lost on restart.
func NewMemorySession() SessionModel {
	return &memorySession{store: make(map[string]memoryEntry)}
}
