//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what []string 을 []pgtype.Text 로 bulk wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgTexts(vs []string) []pgtype.Text {
	if vs == nil {
		return nil
	}
	out := make([]pgtype.Text, len(vs))
	for i, v := range vs {
		out[i] = pgtype.Text{String: v, Valid: true}
	}
	return out
}
