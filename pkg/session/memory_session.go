//ff:type feature=pkg-session type=model
//ff:what 인메모리 세션 구조체 — RWMutex 기반 동시성 제어
package session

import "sync"

type memorySession struct {
	mu    sync.RWMutex
	store map[string]memoryEntry
}
