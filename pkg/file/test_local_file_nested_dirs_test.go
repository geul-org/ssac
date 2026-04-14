//ff:func feature=pkg-file type=test control=sequence
//ff:what 중첩 디렉토리 경로에 파일 Upload/Download가 동작하는지 검증한다
package file

import (
	"bytes"
	"context"
	"io"
	"testing"
)

func TestLocalFile_NestedDirs(t *testing.T) {
	root := t.TempDir()
	f := NewLocalFile(root)
	ctx := context.Background()

	if err := f.Upload(ctx, "a/b/c/d.txt", bytes.NewReader([]byte("nested"))); err != nil {
		t.Fatal(err)
	}

	rc, err := f.Download(ctx, "a/b/c/d.txt")
	if err != nil {
		t.Fatal(err)
	}
	got, _ := io.ReadAll(rc)
	rc.Close()
	if string(got) != "nested" {
		t.Errorf("expected %q, got %q", "nested", string(got))
	}
}
