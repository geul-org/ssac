//ff:func feature=pkg-pgtypex type=util control=sequence
//ff:what pgtype.Interval 의 (Months, Days, Microseconds) 를 ISO 8601 duration 표기로 직렬화한다
package pgtypex

import (
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

// formatIntervalISO renders a pgtype.Interval as ISO 8601 duration. PG's
// month/day/microsecond decomposition is preserved (months cannot be
// expressed in microseconds without a calendar context). Zero intervals
// render as "PT0S" — the standard short form for an empty duration.
func formatIntervalISO(v pgtype.Interval) string {
	if v.Months == 0 && v.Days == 0 && v.Microseconds == 0 {
		return "PT0S"
	}
	var b strings.Builder
	b.WriteByte('P')
	years := v.Months / 12
	months := v.Months % 12
	if years != 0 {
		b.WriteString(strconv.FormatInt(int64(years), 10))
		b.WriteByte('Y')
	}
	if months != 0 {
		b.WriteString(strconv.FormatInt(int64(months), 10))
		b.WriteByte('M')
	}
	if v.Days != 0 {
		b.WriteString(strconv.FormatInt(int64(v.Days), 10))
		b.WriteByte('D')
	}
	if v.Microseconds != 0 {
		b.WriteByte('T')
		us := v.Microseconds
		neg := us < 0
		if neg {
			us = -us
		}
		const usPerHour = int64(3600) * 1_000_000
		const usPerMin = int64(60) * 1_000_000
		hours := us / usPerHour
		us -= hours * usPerHour
		minutes := us / usPerMin
		us -= minutes * usPerMin
		seconds := us / 1_000_000
		fracUs := us - seconds*1_000_000
		if hours != 0 {
			if neg {
				b.WriteByte('-')
			}
			b.WriteString(strconv.FormatInt(hours, 10))
			b.WriteByte('H')
		}
		if minutes != 0 {
			if neg {
				b.WriteByte('-')
			}
			b.WriteString(strconv.FormatInt(minutes, 10))
			b.WriteByte('M')
		}
		if seconds != 0 || fracUs != 0 || (hours == 0 && minutes == 0) {
			if neg {
				b.WriteByte('-')
			}
			b.WriteString(strconv.FormatInt(seconds, 10))
			if fracUs != 0 {
				b.WriteString(formatIntervalFrac(fracUs))
			}
			b.WriteByte('S')
		}
	}
	return b.String()
}
