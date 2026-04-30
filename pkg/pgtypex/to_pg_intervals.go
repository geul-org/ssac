//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what []string 을 []pgtype.Interval 로 bulk wrap 한다 (parse 실패는 해당 슬롯만 Valid:false)
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgIntervals(vs []string) []pgtype.Interval {
	if vs == nil {
		return nil
	}
	out := make([]pgtype.Interval, len(vs))
	for i, v := range vs {
		out[i] = ToPgInterval(v)
	}
	return out
}
