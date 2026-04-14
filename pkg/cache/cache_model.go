//ff:type feature=pkg-cache type=model
//ff:what 캐시 모델 인터페이스 — TTL 기반 키-값 저장소 계약
package cache

import (
	"context"
	"time"
)

// CacheModel provides key-value + TTL storage for data efficiency (caching).
type CacheModel interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}
