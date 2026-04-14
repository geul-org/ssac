//ff:func feature=pkg-session type=util control=sequence
//ff:what 세션에서 key로 value를 조회한다
package session

import "context"

// @func get
// @description 세션에서 key로 value를 조회한다
// @error 404

func Get(req GetRequest) (GetResponse, error) {
	value, err := defaultModel.Get(context.Background(), req.Key)
	return GetResponse{Value: value}, err
}
