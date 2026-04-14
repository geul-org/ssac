//ff:func feature=pkg-file type=util control=sequence
//ff:what 파일을 저장소에 업로드한다
package file

import (
	"context"
	"strings"
)

// @func upload
// @description 파일을 저장소에 업로드한다

func Upload(req UploadRequest) (UploadResponse, error) {
	err := defaultModel.Upload(context.Background(), req.Key, strings.NewReader(req.Body))
	return UploadResponse{Key: req.Key}, err
}
