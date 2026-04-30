//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable *string 을 pgtype.Interval 로 wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgIntervalPtr(v *string) pgtype.Interval {
	if v == nil {
		return pgtype.Interval{}
	}
	out, err := parseIntervalISO(*v)
	if err != nil {
		return pgtype.Interval{}
	}
	return out
}
