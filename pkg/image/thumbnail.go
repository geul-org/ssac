//ff:func feature=pkg-image type=util control=sequence
//ff:what 이미지를 200x200 정사각형으로 크롭하여 PNG 썸네일을 생성한다
package image

import (
	"bytes"

	"github.com/disintegration/imaging"
)

// @func thumbnail
// @description 이미지를 200x200 정사각형으로 크롭하여 PNG 썸네일을 생성한다

func Thumbnail(req ThumbnailRequest) (ThumbnailResponse, error) {
	src, err := imaging.Decode(bytes.NewReader(req.Data))
	if err != nil {
		return ThumbnailResponse{}, err
	}
	thumb := imaging.Fill(src, 200, 200, imaging.Center, imaging.Lanczos)
	var buf bytes.Buffer
	if err := imaging.Encode(&buf, thumb, imaging.PNG); err != nil {
		return ThumbnailResponse{}, err
	}
	return ThumbnailResponse{Data: buf.Bytes()}, nil
}
