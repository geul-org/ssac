//ff:func feature=pkg-pgtypex type=util control=sequence
//ff:what 마이크로초 잔여 (0 < x < 1_000_000) 를 ISO 8601 초 소수부 ".NNNNNN" 로 (trailing zero 제거) 직렬화한다
package pgtypex

import (
	"strconv"
	"strings"
)

func formatIntervalFrac(fracUs int64) string {
	frac := strconv.FormatInt(fracUs, 10)
	pad := 6 - len(frac)
	if pad < 0 {
		pad = 0
	}
	return "." + strings.Repeat("0", pad) + strings.TrimRight(frac, "0")
}
