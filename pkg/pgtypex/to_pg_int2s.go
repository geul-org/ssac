//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what []int16 을 []pgtype.Int2 로 bulk wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgInt2s(vs []int16) []pgtype.Int2 {
	if vs == nil {
		return nil
	}
	out := make([]pgtype.Int2, len(vs))
	for i, v := range vs {
		out[i] = pgtype.Int2{Int16: v, Valid: true}
	}
	return out
}
