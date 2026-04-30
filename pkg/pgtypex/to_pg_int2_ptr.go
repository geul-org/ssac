//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable *int16 을 pgtype.Int2 로 wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgInt2Ptr(v *int16) pgtype.Int2 {
	if v == nil {
		return pgtype.Int2{}
	}
	return pgtype.Int2{Int16: *v, Valid: true}
}
