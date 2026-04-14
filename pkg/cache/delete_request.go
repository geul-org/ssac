//ff:type feature=pkg-cache type=model
//ff:what 캐시 삭제 요청 모델
package cache

type DeleteRequest struct {
	Key string
}
