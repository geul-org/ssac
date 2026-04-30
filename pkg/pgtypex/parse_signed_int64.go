//ff:func feature=pkg-pgtypex type=util control=sequence
//ff:what 부호 포함 십진수 문자열을 int64 로 파싱한다 (interval H/Min 컴포넌트용)
package pgtypex

import "strconv"

func parseSignedInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
