//ff:type feature=pkg-file type=model
//ff:what 파일 모델 인터페이스 — 파일/객체 저장소 계약
package file

import (
	"context"
	"io"
)

// FileModel provides file/object storage operations.
type FileModel interface {
	Upload(ctx context.Context, key string, body io.Reader) error
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}
