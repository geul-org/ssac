//ff:func feature=pkg-auth type=test control=sequence topic=auth-cookie
//ff:what SetAuthCookies / ClearAuthCookies / Extract 속성 검증 (Phase020)
package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// newRecorderCtx builds a minimal gin.Context bound to an httptest recorder
// so SetAuthCookies / ClearAuthCookies can emit Set-Cookie headers we can
// inspect. Kept inline in the test package — shared helper would add a file
// with negligible reuse value given only these three tests consume it.
func newRecorderCtx(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c, w
}

func TestSetAuthCookies_CookieMode_EmitsHostPrefix(t *testing.T) {
	Configure(Config{
		SecretEnv:  "JWT_SECRET",
		AccessTTL:  15 * time.Minute,
		RefreshTTL: 168 * time.Hour,
		Mode:       "cookie",
	})
	t.Cleanup(func() { Configure(Config{}) })

	c, w := newRecorderCtx(t)
	SetAuthCookies(c, "access-jwt", "refresh-jwt")

	cookies := w.Result().Cookies()
	if len(cookies) != 2 {
		t.Fatalf("expected 2 cookies, got %d: %+v", len(cookies), cookies)
	}

	var access, refresh *http.Cookie
	for _, ck := range cookies {
		switch ck.Name {
		case "__Host-access_token":
			access = ck
		case "__Host-refresh_token":
			refresh = ck
		}
	}
	if access == nil {
		t.Fatalf("missing __Host-access_token cookie: %+v", cookies)
	}
	if refresh == nil {
		t.Fatalf("missing __Host-refresh_token cookie: %+v", cookies)
	}

	// Access cookie invariants.
	if !access.HttpOnly {
		t.Errorf("access cookie: HttpOnly must be true")
	}
	if !access.Secure {
		t.Errorf("access cookie: Secure must be true (required by __Host- prefix)")
	}
	if access.Path != "/" {
		t.Errorf("access cookie: Path=%q, want \"/\"", access.Path)
	}
	if access.Domain != "" {
		t.Errorf("access cookie: Domain=%q, want empty (host-only required by __Host-)", access.Domain)
	}
	if access.SameSite != http.SameSiteLaxMode {
		t.Errorf("access cookie: SameSite=%v, want Lax", access.SameSite)
	}
	if access.MaxAge != int((15 * time.Minute).Seconds()) {
		t.Errorf("access cookie: MaxAge=%d, want %d", access.MaxAge, int((15 * time.Minute).Seconds()))
	}

	// Refresh cookie scope.
	if refresh.Path != refreshCookiePath {
		t.Errorf("refresh cookie: Path=%q, want %q", refresh.Path, refreshCookiePath)
	}
	if !refresh.HttpOnly || !refresh.Secure {
		t.Errorf("refresh cookie: HttpOnly/Secure must both be true")
	}
}

func TestSetAuthCookies_BearerMode_NoOp(t *testing.T) {
	Configure(Config{
		SecretEnv:  "JWT_SECRET",
		AccessTTL:  15 * time.Minute,
		RefreshTTL: 168 * time.Hour,
		Mode:       "bearer",
	})
	t.Cleanup(func() { Configure(Config{}) })

	c, w := newRecorderCtx(t)
	SetAuthCookies(c, "access", "refresh")

	if cookies := w.Result().Cookies(); len(cookies) != 0 {
		t.Fatalf("bearer mode must be a no-op, got %d cookies: %+v", len(cookies), cookies)
	}
	// Also verify no Set-Cookie header was written at all (belt + suspenders
	// — Cookies() filters malformed headers).
	if got := w.Header().Get("Set-Cookie"); got != "" {
		t.Fatalf("bearer mode emitted Set-Cookie header: %q", got)
	}
}

func TestSetAuthCookies_UnsetMode_NoOp(t *testing.T) {
	// Pre-Phase020 Configure (no Mode set) must not accidentally emit
	// cookies. The generated boot code fills Mode explicitly; this test
	// guards against future reflows that forget to set it.
	Configure(Config{SecretEnv: "JWT_SECRET", AccessTTL: time.Minute, RefreshTTL: time.Hour})
	t.Cleanup(func() { Configure(Config{}) })

	c, w := newRecorderCtx(t)
	SetAuthCookies(c, "a", "r")
	if cookies := w.Result().Cookies(); len(cookies) != 0 {
		t.Fatalf("unset mode must be a no-op, got %d cookies", len(cookies))
	}
}

func TestClearAuthCookies_CookieMode_MaxAgeNegative(t *testing.T) {
	Configure(Config{
		SecretEnv:  "JWT_SECRET",
		AccessTTL:  15 * time.Minute,
		RefreshTTL: 168 * time.Hour,
		Mode:       "cookie",
	})
	t.Cleanup(func() { Configure(Config{}) })

	c, w := newRecorderCtx(t)
	ClearAuthCookies(c)

	// net/http renders MaxAge<0 as "Max-Age=0; Expires=<epoch>". We check
	// the presence of the MaxAge=0 token on every Set-Cookie header line
	// because cookies parsed back via w.Result().Cookies() drop MaxAge=0
	// (browsers treat it as a delete signal, which is exactly what we
	// want to assert).
	for _, line := range w.Header().Values("Set-Cookie") {
		if !strings.Contains(line, "Max-Age=0") {
			t.Errorf("expected Max-Age=0 delete marker in %q", line)
		}
	}

	// Still two headers (one per cookie name) and each carries the same
	// Path that was used for issuance — browsers require Path match to
	// delete.
	headers := w.Header().Values("Set-Cookie")
	if len(headers) != 2 {
		t.Fatalf("expected 2 Set-Cookie headers, got %d: %v", len(headers), headers)
	}
	foundAccessPath := false
	foundRefreshPath := false
	for _, line := range headers {
		if strings.Contains(line, "__Host-access_token=") && strings.Contains(line, "Path=/;") {
			foundAccessPath = true
		}
		if strings.Contains(line, "__Host-refresh_token=") && strings.Contains(line, "Path=/auth/refresh") {
			foundRefreshPath = true
		}
	}
	if !foundAccessPath {
		t.Errorf("access delete cookie missing Path=/: %v", headers)
	}
	if !foundRefreshPath {
		t.Errorf("refresh delete cookie missing Path=/auth/refresh: %v", headers)
	}
}

func TestClearAuthCookies_BearerMode_NoOp(t *testing.T) {
	Configure(Config{Mode: "bearer"})
	t.Cleanup(func() { Configure(Config{}) })

	c, w := newRecorderCtx(t)
	ClearAuthCookies(c)
	if got := w.Header().Get("Set-Cookie"); got != "" {
		t.Fatalf("bearer mode ClearAuthCookies must be no-op, got %q", got)
	}
}

func TestExtractAccessFromCookie_RoundTrip(t *testing.T) {
	Configure(Config{Mode: "cookie", AccessTTL: time.Minute, RefreshTTL: time.Hour})
	t.Cleanup(func() { Configure(Config{}) })

	c, _ := newRecorderCtx(t)
	// URL-escape the value the same way gin.SetCookie does — tests go
	// through the full cookie jar, so we build the request cookie by hand.
	c.Request.AddCookie(&http.Cookie{Name: "__Host-access_token", Value: "hello"})
	if got := ExtractAccessFromCookie(c); got != "hello" {
		t.Errorf("ExtractAccessFromCookie=%q, want %q", got, "hello")
	}
	if got := ExtractRefreshFromCookie(c); got != "" {
		t.Errorf("ExtractRefreshFromCookie should be empty (no cookie set), got %q", got)
	}
}

func TestResolvedCookieAttrs_FillsDefaults(t *testing.T) {
	cfg := Config{AccessTTL: 5 * time.Minute, RefreshTTL: time.Hour}
	a := resolvedCookieAttrs(cfg)
	if a.AccessName != "__Host-access_token" {
		t.Errorf("AccessName default=%q", a.AccessName)
	}
	if a.RefreshName != "__Host-refresh_token" {
		t.Errorf("RefreshName default=%q", a.RefreshName)
	}
	if a.SameSite != http.SameSiteLaxMode {
		t.Errorf("SameSite default=%v", a.SameSite)
	}
	if a.AccessTTL != 5*time.Minute {
		t.Errorf("AccessTTL fallback=%v", a.AccessTTL)
	}
	if a.RefreshTTL != time.Hour {
		t.Errorf("RefreshTTL fallback=%v", a.RefreshTTL)
	}
}
