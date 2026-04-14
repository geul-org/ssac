//ff:type feature=pkg-cache type=model
//ff:what 인메모리 캐시 구조체 — RWMutex 기반 동시성 제어
package cache

import "sync"

type memoryCache struct {
	mu    sync.RWMutex
	store map[string]memoryEntry
}
