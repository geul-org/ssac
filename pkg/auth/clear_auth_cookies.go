//ff:func feature=pkg-auth type=util control=sequence topic=auth-cookie
//ff:what ClearAuthCookies — Logout 시 access/refresh 쿠키 Max-Age=-1 로 제거 (Phase020)
package auth

import (
	"github.com/gin-gonic/gin"
)

// ClearAuthCookies removes the access + refresh cookies by emitting
// Set-Cookie headers with MaxAge=-1 (RFC 6265 §5.3 — "delete the cookie").
// Called automatically by the generated logout handler when Mode != "bearer".
//
// Cookie removal requires the Path attribute to match the original issuance
// (browsers key cookies on name + path + domain). We therefore emit:
//
//   - name = __Host-access_token,  path = "/"
//   - name = __Host-refresh_token, path = "/auth/refresh"
//
// MaxAge=-1 is preferred over MaxAge=0 because Go's net/http renders it as
// the explicit `Max-Age=0` + `Expires=<epoch>` combo, which older browsers
// interpret unambiguously as an expiration signal.
func ClearAuthCookies(c *gin.Context) {
	cfg := currentConfig()
	if cfg.Mode == "bearer" || cfg.Mode == "" {
		return
	}
	attrs := resolvedCookieAttrs(cfg)

	c.SetSameSite(attrs.SameSite)
	c.SetCookie(attrs.AccessName, "", -1, "/", "", true, true)

	c.SetSameSite(attrs.SameSite)
	c.SetCookie(attrs.RefreshName, "", -1, refreshCookiePath, "", true, true)
}
