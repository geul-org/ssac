//ff:type feature=pkg-cache type=model
//ff:what 메모리 캐시 항목 — 값과 만료 시각 보관
package cache

import "time"

type memoryEntry struct {
	value     string
	expiresAt time.Time
}
