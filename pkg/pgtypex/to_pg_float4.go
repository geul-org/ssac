//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL float32 를 pgtype.Float4 (Valid:true) 로 wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgFloat4(v float32) pgtype.Float4 {
	return pgtype.Float4{Float32: v, Valid: true}
}
