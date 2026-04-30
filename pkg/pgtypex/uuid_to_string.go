//ff:func feature=pkg-pgtypex type=util control=sequence
//ff:what pgtype.UUID 를 canonical 8-4-4-4-12 hex 문자열로 직렬화한다 (NULL → "")
package pgtypex

import (
	"encoding/hex"

	"github.com/jackc/pgx/v5/pgtype"
)

// UUIDToString returns the canonical 8-4-4-4-12 form for a Valid pgtype.UUID
// or "" for SQL NULL. yongol SSaC emit calls it at OPA Owners-map sites
// where the resource id must be string-serialised regardless of its
// underlying PG type.
func UUIDToString(v pgtype.UUID) string {
	if !v.Valid {
		return ""
	}
	b := v.Bytes
	dst := make([]byte, 36)
	hex.Encode(dst[0:8], b[0:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], b[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], b[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], b[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:36], b[10:16])
	return string(dst)
}
