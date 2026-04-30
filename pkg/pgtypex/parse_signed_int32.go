//ff:func feature=pkg-pgtypex type=util control=sequence
//ff:what 부호 포함 십진수 문자열을 int32 로 파싱한다 (interval Y/M/D 컴포넌트용)
package pgtypex

import "strconv"

func parseSignedInt32(s string) (int32, error) {
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}
