//ff:func feature=pkg-cache type=util control=sequence
//ff:what 캐시에 key-value를 저장한다
package cache

import (
	"context"
	"time"
)

// @func set
// @description 캐시에 key-value를 저장한다

func Set(ctx context.Context, req SetRequest) (SetResponse, error) {
	return SetResponse{}, defaultModel.Set(ctx, req.Key, req.Value, time.Duration(req.TTL)*time.Second)
}
