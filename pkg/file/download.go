//ff:func feature=pkg-file type=util control=sequence
//ff:what 저장소에서 파일을 다운로드한다
package file

import (
	"context"
	"io"
)

// @func download
// @description 저장소에서 파일을 다운로드한다
// @error 404

func Download(ctx context.Context, req DownloadRequest) (DownloadResponse, error) {
	rc, err := defaultModel.Download(ctx, req.Key)
	if err != nil {
		return DownloadResponse{}, err
	}
	defer rc.Close()
	data, err := io.ReadAll(rc)
	return DownloadResponse{Body: string(data)}, err
}
