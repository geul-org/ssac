//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL int16 을 pgtype.Int2 (Valid:true) 로 wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgInt2(v int16) pgtype.Int2 {
	return pgtype.Int2{Int16: v, Valid: true}
}
