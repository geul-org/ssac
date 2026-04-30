//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL ISO 8601 duration string 을 pgtype.Interval (Valid:true) 로 wrap 한다 (parse 실패는 Valid:false)
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgInterval(v string) pgtype.Interval {
	out, err := parseIntervalISO(v)
	if err != nil {
		return pgtype.Interval{}
	}
	return out
}
