//ff:func feature=pkg-text type=util control=sequence
//ff:what HTML에서 위험한 태그와 속성을 제거한다 (XSS 방지)
package text

import "github.com/microcosm-cc/bluemonday"

// @func sanitizeHTML
// @description HTML에서 위험한 태그와 속성을 제거한다 (XSS 방지)

func SanitizeHTML(req SanitizeHTMLRequest) (SanitizeHTMLResponse, error) {
	p := bluemonday.UGCPolicy()
	return SanitizeHTMLResponse{Sanitized: p.Sanitize(req.HTML)}, nil
}
