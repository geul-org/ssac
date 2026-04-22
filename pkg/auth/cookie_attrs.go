//ff:type feature=pkg-auth type=model topic=auth-cookie
//ff:what Auth 쿠키 발급 속성 (Phase020) — AccessName / RefreshName / SameSite / TTL
package auth

import (
	"net/http"
	"time"
)

// CookieAttrs bundles the session-cookie parameters applied when the backend
// runs in `cookie` or `hybrid` auth mode (Phase020). It is injected through
// Config.CookieAttrs on Configure(); SetAuthCookies / ClearAuthCookies read
// the current snapshot on every call so tests can mutate attrs via Configure
// without lifetime-managed state.
//
// Defaults (applied in Configure if the struct is zero-valued):
//   - AccessName:  "__Host-access_token"
//   - RefreshName: "__Host-refresh_token"
//   - SameSite:    http.SameSiteLaxMode
//   - AccessTTL / RefreshTTL: fall back to Config.AccessTTL / RefreshTTL.
//
// The `__Host-` prefix is intentional — it forces the browser to reject the
// cookie unless Secure=true, Path=/ (for access), and Domain is absent. This
// turns operator mis-configuration (HTTP dev URL, Domain typo) into a
// browser-level failure instead of a silent security downgrade.
type CookieAttrs struct {
	// AccessName is the Set-Cookie name for the access token. Default
	// "__Host-access_token". Callers may override, but names that don't
	// start with "__Host-" lose the prefix-level browser enforcement.
	AccessName string
	// RefreshName is the Set-Cookie name for the refresh token. Default
	// "__Host-refresh_token".
	RefreshName string
	// SameSite controls the Set-Cookie SameSite attribute. Default Lax
	// (safe for top-level navigation, blocks cross-site POST). Strict
	// breaks cross-site navigation logins; None requires Secure and
	// weakens CSRF protection (use only for cross-origin SPA).
	SameSite http.SameSite
	// AccessTTL overrides the Max-Age of the access cookie. Zero means
	// "fall back to Config.AccessTTL".
	AccessTTL time.Duration
	// RefreshTTL overrides the Max-Age of the refresh cookie. Zero means
	// "fall back to Config.RefreshTTL".
	RefreshTTL time.Duration
}

// resolvedCookieAttrs fills zero-valued fields with defaults. Internal
// helper shared by SetAuthCookies / ClearAuthCookies / ExtractAccessFromCookie
// so every caller sees the same name even when Configure was never called.
func resolvedCookieAttrs(cfg Config) CookieAttrs {
	a := cfg.CookieAttrs
	if a.AccessName == "" {
		a.AccessName = "__Host-access_token"
	}
	if a.RefreshName == "" {
		a.RefreshName = "__Host-refresh_token"
	}
	if a.SameSite == 0 {
		a.SameSite = http.SameSiteLaxMode
	}
	if a.AccessTTL == 0 {
		a.AccessTTL = cfg.AccessTTL
	}
	if a.RefreshTTL == 0 {
		a.RefreshTTL = cfg.RefreshTTL
	}
	return a
}
