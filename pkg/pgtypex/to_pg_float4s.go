//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what []float32 를 []pgtype.Float4 로 bulk wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgFloat4s(vs []float32) []pgtype.Float4 {
	if vs == nil {
		return nil
	}
	out := make([]pgtype.Float4, len(vs))
	for i, v := range vs {
		out[i] = pgtype.Float4{Float32: v, Valid: true}
	}
	return out
}
