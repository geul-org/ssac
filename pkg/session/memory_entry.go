//ff:type feature=pkg-session type=model
//ff:what 메모리 세션 항목 — 값과 만료 시각 보관
package session

import "time"

type memoryEntry struct {
	value     string
	expiresAt time.Time
}
