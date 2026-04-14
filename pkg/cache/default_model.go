//ff:func feature=pkg-cache type=util control=sequence
//ff:what @call Func용 기본 캐시 모델 — Init으로 주입
package cache

var defaultModel CacheModel

// Init sets the package-level CacheModel used by @call Func wrappers.
func Init(model CacheModel) {
	defaultModel = model
}
