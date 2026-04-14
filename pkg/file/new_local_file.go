//ff:func feature=pkg-file type=loader control=sequence
//ff:what 로컬 파일 저장소 생성 — 루트 디렉토리 경로로 인스턴스 반환
package file

// NewLocalFile creates a FileModel backed by the local filesystem.
func NewLocalFile(root string) FileModel {
	return &localFile{root: root}
}
