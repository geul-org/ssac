//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what []string 을 []pgtype.Numeric 로 bulk wrap 한다 (parse 실패는 해당 슬롯만 Valid:false)
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgNumerics(vs []string) []pgtype.Numeric {
	if vs == nil {
		return nil
	}
	out := make([]pgtype.Numeric, len(vs))
	for i, v := range vs {
		out[i] = ToPgNumeric(v)
	}
	return out
}
