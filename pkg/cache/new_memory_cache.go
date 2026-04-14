//ff:func feature=pkg-cache type=loader control=sequence
//ff:what 인메모리 캐시 생성 — 재시작 시 데이터 소멸
package cache

// NewMemoryCache creates an in-memory CacheModel. Data is lost on restart.
func NewMemoryCache() CacheModel {
	return &memoryCache{store: make(map[string]memoryEntry)}
}
