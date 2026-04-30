//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable *float64 를 pgtype.Float8 로 wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgFloat8Ptr(v *float64) pgtype.Float8 {
	if v == nil {
		return pgtype.Float8{}
	}
	return pgtype.Float8{Float64: *v, Valid: true}
}
