//ff:func feature=pkg-session type=util control=sequence
//ff:what 세션에 key-value를 저장한다
package session

import (
	"context"
	"time"
)

// @func set
// @description 세션에 key-value를 저장한다

func Set(req SetRequest) (SetResponse, error) {
	return SetResponse{}, defaultModel.Set(context.Background(), req.Key, req.Value, time.Duration(req.TTL)*time.Second)
}
