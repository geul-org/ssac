//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what []int32 를 []pgtype.Int4 로 bulk wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgInt4s(vs []int32) []pgtype.Int4 {
	if vs == nil {
		return nil
	}
	out := make([]pgtype.Int4, len(vs))
	for i, v := range vs {
		out[i] = pgtype.Int4{Int32: v, Valid: true}
	}
	return out
}
