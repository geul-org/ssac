//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what []bool 을 []pgtype.Bool 로 bulk wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgBools(vs []bool) []pgtype.Bool {
	if vs == nil {
		return nil
	}
	out := make([]pgtype.Bool, len(vs))
	for i, v := range vs {
		out[i] = pgtype.Bool{Bool: v, Valid: true}
	}
	return out
}
