//ff:func feature=pkg-file type=util control=sequence
//ff:what 로컬 Delete — 파일 삭제
package file

import (
	"context"
	"os"
	"path/filepath"
)

func (f *localFile) Delete(_ context.Context, key string) error {
	path := filepath.Join(f.root, key)
	return os.Remove(path)
}
