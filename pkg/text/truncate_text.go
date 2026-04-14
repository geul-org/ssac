//ff:func feature=pkg-text type=util control=sequence
//ff:what 유니코드 안전하게 텍스트를 자른다
package text

// @func truncateText
// @description 유니코드 안전하게 텍스트를 자른다

func TruncateText(req TruncateTextRequest) (TruncateTextResponse, error) {
	suffix := req.Suffix
	if suffix == "" {
		suffix = "..."
	}
	runes := []rune(req.Text)
	if len(runes) <= req.MaxLength {
		return TruncateTextResponse{Truncated: req.Text}, nil
	}
	return TruncateTextResponse{Truncated: string(runes[:req.MaxLength]) + suffix}, nil
}
