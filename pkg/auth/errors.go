//ff:type feature=pkg-auth type=data topic=auth-refresh
//ff:what auth 패키지 sentinel errors (store 구현체 공용)
package auth

import "errors"

// ErrEmptyMatcher is returned by RefreshStore.RevokeAll when the matcher is
// empty. Prevents accidental full-table revocation.
var ErrEmptyMatcher = errors.New("refresh store: empty matcher rejected")
