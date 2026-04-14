//ff:func feature=pkg-text type=util control=sequence
//ff:what 텍스트를 URL-safe slug로 변환한다
package text

import "github.com/gosimple/slug"

// @func generateSlug
// @description 텍스트를 URL-safe slug로 변환한다

func GenerateSlug(req GenerateSlugRequest) (GenerateSlugResponse, error) {
	return GenerateSlugResponse{Slug: slug.Make(req.Text)}, nil
}
