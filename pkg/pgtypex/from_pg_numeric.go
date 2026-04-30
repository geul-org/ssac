//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL pgtype.Numeric 를 textual string 으로 unwrap 한다 (NaN → "NaN")
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

// FromPgNumeric returns the canonical textual form of a Valid pgtype.Numeric.
// MarshalJSON is reused because it is the only public path that emits the
// authoritative decimal text without depending on driver-side encoding.
// NaN renders as the unquoted token "NaN" matching MarshalJSON's quoted
// "NaN" minus the surrounding quotes.
func FromPgNumeric(v pgtype.Numeric) string {
	if !v.Valid {
		return ""
	}
	buf, err := v.MarshalJSON()
	if err != nil {
		return ""
	}
	s := string(buf)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
