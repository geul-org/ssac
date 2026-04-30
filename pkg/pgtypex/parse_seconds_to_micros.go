//ff:func feature=pkg-pgtypex type=util control=sequence
//ff:what ISO 8601 duration 초 컴포넌트 (소수점 6자리까지) 를 마이크로초 정수로 변환한다
package pgtypex

import (
	"strconv"
	"strings"
)

func parseSecondsToMicros(s string) (int64, error) {
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	} else if strings.HasPrefix(s, "+") {
		s = s[1:]
	}
	intPart := s
	fracPart := ""
	if i := strings.IndexByte(s, '.'); i >= 0 {
		intPart = s[:i]
		fracPart = s[i+1:]
	}
	whole, err := strconv.ParseInt(intPart, 10, 64)
	if err != nil {
		return 0, err
	}
	micros := whole * 1_000_000
	if fracPart != "" {
		if len(fracPart) > 6 {
			fracPart = fracPart[:6]
		}
		for len(fracPart) < 6 {
			fracPart += "0"
		}
		f, err := strconv.ParseInt(fracPart, 10, 64)
		if err != nil {
			return 0, err
		}
		micros += f
	}
	if neg {
		micros = -micros
	}
	return micros, nil
}
