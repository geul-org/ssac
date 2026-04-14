//ff:func feature=pkg-file type=test control=sequence
//ff:what 존재하지 않는 파일 Download 시 에러를 반환하는지 검증한다
package file

import (
	"context"
	"testing"
)

func TestLocalFile_DownloadNotFound(t *testing.T) {
	root := t.TempDir()
	f := NewLocalFile(root)
	ctx := context.Background()

	_, err := f.Download(ctx, "nonexistent.txt")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
