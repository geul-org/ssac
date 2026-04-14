//ff:type feature=pkg-cache type=model
//ff:what 캐시 저장 요청 모델
package cache

type SetRequest struct {
	Key   string
	Value string
	TTL   int64
}
