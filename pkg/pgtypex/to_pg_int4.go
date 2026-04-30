//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL int32 를 pgtype.Int4 (Valid:true) 로 wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgInt4(v int32) pgtype.Int4 {
	return pgtype.Int4{Int32: v, Valid: true}
}
