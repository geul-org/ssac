//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL float64 를 pgtype.Float8 (Valid:true) 로 wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgFloat8(v float64) pgtype.Float8 {
	return pgtype.Float8{Float64: v, Valid: true}
}
