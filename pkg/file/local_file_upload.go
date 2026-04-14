//ff:func feature=pkg-file type=util control=selection
//ff:what 로컬 Upload — 디렉토리 생성 후 파일 쓰기, 오류 시 분기
package file

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

func (f *localFile) Upload(_ context.Context, key string, body io.Reader) error {
	path := filepath.Join(f.root, key)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, body)
	return err
}
