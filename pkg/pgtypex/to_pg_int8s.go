//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what []int64 를 []pgtype.Int8 로 bulk wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgInt8s(vs []int64) []pgtype.Int8 {
	if vs == nil {
		return nil
	}
	out := make([]pgtype.Int8, len(vs))
	for i, v := range vs {
		out[i] = pgtype.Int8{Int64: v, Valid: true}
	}
	return out
}
