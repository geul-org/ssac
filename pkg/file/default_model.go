//ff:func feature=pkg-file type=util control=sequence
//ff:what @call Func용 기본 파일 모델 — Init으로 주입
package file

var defaultModel FileModel

// Init sets the package-level FileModel used by @call Func wrappers.
func Init(model FileModel) {
	defaultModel = model
}
