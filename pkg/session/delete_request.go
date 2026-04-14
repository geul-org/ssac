//ff:type feature=pkg-session type=model
//ff:what 세션 삭제 요청 모델
package session

type DeleteRequest struct {
	Key string
}
