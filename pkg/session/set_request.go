//ff:type feature=pkg-session type=model
//ff:what 세션 저장 요청 모델
package session

type SetRequest struct {
	Key   string
	Value string
	TTL   int64
}
