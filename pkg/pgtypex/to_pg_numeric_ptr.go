//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable *string 을 pgtype.Numeric 로 wrap 한다 (nil → Valid:false)
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgNumericPtr(v *string) pgtype.Numeric {
	if v == nil {
		return pgtype.Numeric{}
	}
	var n pgtype.Numeric
	if err := n.Scan(*v); err != nil {
		return pgtype.Numeric{}
	}
	return n
}
