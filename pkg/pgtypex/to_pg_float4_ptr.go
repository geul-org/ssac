//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable *float32 를 pgtype.Float4 로 wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgFloat4Ptr(v *float32) pgtype.Float4 {
	if v == nil {
		return pgtype.Float4{}
	}
	return pgtype.Float4{Float32: *v, Valid: true}
}
