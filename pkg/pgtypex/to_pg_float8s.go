//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what []float64 를 []pgtype.Float8 로 bulk wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgFloat8s(vs []float64) []pgtype.Float8 {
	if vs == nil {
		return nil
	}
	out := make([]pgtype.Float8, len(vs))
	for i, v := range vs {
		out[i] = pgtype.Float8{Float64: v, Valid: true}
	}
	return out
}
