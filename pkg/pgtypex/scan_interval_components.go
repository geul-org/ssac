//ff:func feature=pkg-pgtypex type=util control=sequence
//ff:what ISO 8601 duration 본문을 단위별로 누적 파싱한다 (date / time half 분리)
package pgtypex

import "errors"

// scanIntervalComponents walks a half (date or time) of an ISO 8601
// duration body, dispatching each numeric+unit pair to the matching
// accumulator pointer. Unknown units / malformed numbers return an error
// caught by parseIntervalISO.
func scanIntervalComponents(s string, isTime bool, months, days *int32, us *int64) error {
	if s == "" {
		return nil
	}
	i := 0
	for i < len(s) {
		j := i
		if s[j] == '-' || s[j] == '+' {
			j++
		}
		for j < len(s) && (s[j] >= '0' && s[j] <= '9' || s[j] == '.') {
			j++
		}
		if j == i || j >= len(s) {
			return errors.New("pgtypex: malformed ISO 8601 duration component")
		}
		num := s[i:j]
		unit := s[j]
		i = j + 1
		switch {
		case !isTime && unit == 'Y':
			n, err := parseSignedInt32(num)
			if err != nil {
				return err
			}
			if months != nil {
				*months += n * 12
			}
		case !isTime && unit == 'M':
			n, err := parseSignedInt32(num)
			if err != nil {
				return err
			}
			if months != nil {
				*months += n
			}
		case !isTime && unit == 'D':
			n, err := parseSignedInt32(num)
			if err != nil {
				return err
			}
			if days != nil {
				*days += n
			}
		case isTime && unit == 'H':
			n, err := parseSignedInt64(num)
			if err != nil {
				return err
			}
			if us != nil {
				*us += n * 3600 * 1_000_000
			}
		case isTime && unit == 'M':
			n, err := parseSignedInt64(num)
			if err != nil {
				return err
			}
			if us != nil {
				*us += n * 60 * 1_000_000
			}
		case isTime && unit == 'S':
			micros, err := parseSecondsToMicros(num)
			if err != nil {
				return err
			}
			if us != nil {
				*us += micros
			}
		default:
			return errors.New("pgtypex: unknown ISO 8601 duration unit '" + string(unit) + "'")
		}
	}
	return nil
}
