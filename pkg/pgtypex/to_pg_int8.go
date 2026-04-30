//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL int64 를 pgtype.Int8 (Valid:true) 로 wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgInt8(v int64) pgtype.Int8 {
	return pgtype.Int8{Int64: v, Valid: true}
}
