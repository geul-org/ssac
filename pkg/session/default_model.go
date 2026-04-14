//ff:func feature=pkg-session type=util control=sequence
//ff:what @call Func용 기본 세션 모델 — Init으로 주입
package session

var defaultModel SessionModel

// Init sets the package-level SessionModel used by @call Func wrappers.
func Init(model SessionModel) {
	defaultModel = model
}
