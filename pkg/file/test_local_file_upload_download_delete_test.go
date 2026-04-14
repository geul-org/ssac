//ff:func feature=pkg-file type=test control=sequence
//ff:what LocalFile의 Upload/Download/Delete 기본 동작을 검증한다
package file

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalFile_UploadDownloadDelete(t *testing.T) {
	root := t.TempDir()
	f := NewLocalFile(root)
	ctx := context.Background()

	content := []byte("hello world")
	if err := f.Upload(ctx, "test/a.txt", bytes.NewReader(content)); err != nil {
		t.Fatal(err)
	}

	// Verify file exists on disk.
	if _, err := os.Stat(filepath.Join(root, "test/a.txt")); err != nil {
		t.Fatalf("uploaded file not found: %v", err)
	}

	rc, err := f.Download(ctx, "test/a.txt")
	if err != nil {
		t.Fatal(err)
	}
	got, _ := io.ReadAll(rc)
	rc.Close()
	if string(got) != "hello world" {
		t.Errorf("expected %q, got %q", "hello world", string(got))
	}

	if err := f.Delete(ctx, "test/a.txt"); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(root, "test/a.txt")); !os.IsNotExist(err) {
		t.Error("file should be deleted")
	}
}
