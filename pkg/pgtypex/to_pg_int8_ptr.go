//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable *int64 를 pgtype.Int8 로 wrap 한다 (nil → Valid:false)
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgInt8Ptr(v *int64) pgtype.Int8 {
	if v == nil {
		return pgtype.Int8{}
	}
	return pgtype.Int8{Int64: *v, Valid: true}
}
