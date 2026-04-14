//ff:type feature=pkg-text type=model
//ff:what 텍스트 자르기 요청 모델
package text

type TruncateTextRequest struct {
	Text      string
	MaxLength int
	Suffix    string // 말줄임 (기본 "...")
}
