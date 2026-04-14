//ff:func feature=pkg-cache type=util control=sequence
//ff:what 캐시에서 key로 value를 조회한다
package cache

import "context"

// @func get
// @description 캐시에서 key로 value를 조회한다

func Get(req GetRequest) (GetResponse, error) {
	value, err := defaultModel.Get(context.Background(), req.Key)
	return GetResponse{Value: value}, err
}
