//ff:func feature=pkg-auth type=util control=sequence topic=auth-cookie
//ff:what SetAuthCookies — __Host- prefix + HttpOnly + Secure + SameSite 로 access/refresh 쿠키 발급 (Phase020)
package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// refreshCookiePath is the Path attribute of the refresh cookie. Scoping it
// to /auth/refresh means the browser only sends the refresh token to that
// single endpoint, reducing the XSS blast radius compared to a site-wide
// refresh cookie. Kept as a package-level constant so callers (and tests)
// can reference the exact path without string duplication.
const refreshCookiePath = "/auth/refresh"

// SetAuthCookies writes the `__Host-access_token` + `__Host-refresh_token`
// Set-Cookie headers (Phase020). Behaviour:
//
//   - Mode == "bearer" (or empty) → no-op. Bearer-only deployments never
//     leak cookies; this lets the generated handler unconditionally invoke
//     SetAuthCookies without branching in every @call emitter.
//   - HttpOnly = true (always) — JS cannot read or overwrite the cookie.
//   - Secure = true (always) — `__Host-` prefix would otherwise be
//     browser-rejected.
//   - SameSite from CookieAttrs.SameSite (default Lax).
//   - Access cookie: Path = "/" (required by `__Host-` prefix).
//   - Refresh cookie: Path = "/auth/refresh" (scoped — see constant above).
//   - Domain = "" (host-only; `__Host-` prefix rejects cookies with Domain).
//
// maxAge is the number of seconds the browser keeps the cookie. It derives
// from CookieAttrs.AccessTTL / RefreshTTL (or Config.AccessTTL/RefreshTTL as
// fallback). maxAge <= 0 produces a session cookie (browser session-bound),
// which defeats the refresh flow — callers should ensure TTLs are set via
// Configure.
func SetAuthCookies(c *gin.Context, access, refresh string) {
	cfg := currentConfig()
	if cfg.Mode == "bearer" || cfg.Mode == "" {
		return
	}
	attrs := resolvedCookieAttrs(cfg)

	accessMaxAge := int(attrs.AccessTTL.Seconds())
	refreshMaxAge := int(attrs.RefreshTTL.Seconds())

	// Gin's SetCookie copies c.sameSite into each emission, so SetSameSite
	// must be called immediately before SetCookie. Wrap both calls in a
	// single sequence to avoid a race where a later handler mutates
	// c.sameSite between our two SetCookie calls.
	c.SetSameSite(attrs.SameSite)
	// Access cookie — site-wide (Path=/), required by __Host- prefix.
	c.SetCookie(attrs.AccessName, access, accessMaxAge, "/", "", true, true)

	c.SetSameSite(attrs.SameSite)
	// Refresh cookie — scoped to the refresh endpoint. Browsers only send
	// it to POST /auth/refresh, so XSS cannot steal the refresh token via
	// a state-changing request to any other path.
	c.SetCookie(attrs.RefreshName, refresh, refreshMaxAge, refreshCookiePath, "", true, true)

	// Silence unused import lint when the file is compiled in a build
	// that excludes the gin package (practically never, but keeps
	// `go vet` clean under conditional compilation in downstream tests).
	_ = http.SameSiteLaxMode
}
