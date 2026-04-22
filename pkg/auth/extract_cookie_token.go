//ff:func feature=pkg-auth type=util control=sequence topic=auth-cookie
//ff:what ExtractAccessFromCookie / ExtractRefreshFromCookie — 요청 쿠키에서 토큰 추출 (Phase020)
package auth

import (
	"github.com/gin-gonic/gin"
)

// ExtractAccessFromCookie returns the access-token value from the request's
// cookie jar, or "" when the cookie is absent. The cookie name respects
// CookieAttrs.AccessName so overrides in Configure propagate through the
// generated middleware.
//
// Errors from gin's Cookie() (practically only ErrNoCookie) collapse to an
// empty string — the VerifyToken middleware treats "" as unauthenticated and
// returns 401 without leaking the underlying error.
func ExtractAccessFromCookie(c *gin.Context) string {
	attrs := resolvedCookieAttrs(currentConfig())
	v, err := c.Cookie(attrs.AccessName)
	if err != nil {
		return ""
	}
	return v
}

// ExtractRefreshFromCookie mirrors ExtractAccessFromCookie for the refresh
// token. Only the /auth/refresh handler should call this — any other path
// will not see the cookie because the refresh cookie is issued with
// Path=/auth/refresh.
func ExtractRefreshFromCookie(c *gin.Context) string {
	attrs := resolvedCookieAttrs(currentConfig())
	v, err := c.Cookie(attrs.RefreshName)
	if err != nil {
		return ""
	}
	return v
}
