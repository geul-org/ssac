//ff:type feature=pkg-session type=model
//ff:what 세션 모델 인터페이스 — TTL 기반 사용자 상태 저장소 계약
package session

import (
	"context"
	"time"
)

// SessionModel provides key-value + TTL storage for user-bound state (login, cart, etc.).
type SessionModel interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}
