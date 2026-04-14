//ff:func feature=pkg-file type=util control=sequence
//ff:what 로컬 Download — 파일 열기
package file

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

func (f *localFile) Download(_ context.Context, key string) (io.ReadCloser, error) {
	path := filepath.Join(f.root, key)
	return os.Open(path)
}
