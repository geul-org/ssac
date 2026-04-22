//ff:type feature=pkg-auth type=model topic=auth-jwt
//ff:what JWT 발급/검증에 필요한 설정 저장
package auth

import (
	"sync"
	"time"
)

// Config captures runtime JWT settings. Only the env variable NAME is stored;
// the actual secret value is read from os.Getenv on every Issue/Refresh/Verify
// call. This preserves the ability to rotate secrets without re-calling
// Configure and keeps test isolation simple via t.Setenv.
type Config struct {
	// SecretEnv is the name of the environment variable holding the signing
	// secret (e.g. "JWT_SECRET"). Empty SecretEnv causes issue/verify calls
	// to return an error.
	SecretEnv string
	// AccessTTL is the lifetime of access tokens issued via IssueToken.
	AccessTTL time.Duration
	// RefreshTTL is the lifetime of refresh tokens issued via RefreshToken.
	RefreshTTL time.Duration
}

var (
	cfgMu     sync.RWMutex
	globalCfg Config
)

// Configure sets the package-level JWT configuration. Intended to be called
// once during application bootstrap (e.g. from main.go). Subsequent calls
// overwrite previous values (last-write-wins); callers needing strict
// one-time-init semantics should enforce that at their boot layer.
func Configure(cfg Config) {
	cfgMu.Lock()
	globalCfg = cfg
	cfgMu.Unlock()
}

// currentConfig returns a snapshot of the package-level Config under a
// read-lock. Internal helper shared by issue/refresh/verify.
func currentConfig() Config {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return globalCfg
}
