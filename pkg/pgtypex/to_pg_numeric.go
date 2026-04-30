//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL string 을 pgtype.Numeric (Valid:true) 로 wrap 한다 (parse 실패는 Valid:false)
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

// ToPgNumeric scans the textual representation through pgtype.Numeric.Scan
// to preserve arbitrary precision. Caller-provided strings that fail to
// parse fall back to the zero (Valid:false) value rather than panicking;
// upstream validate D-12 plus OpenAPI `format: number` guard against
// malformed input before the bridge is reached.
func ToPgNumeric(v string) pgtype.Numeric {
	var n pgtype.Numeric
	if err := n.Scan(v); err != nil {
		return pgtype.Numeric{}
	}
	return n
}
