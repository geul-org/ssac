//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL pgtype.Interval 을 ISO 8601 duration string 으로 unwrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func FromPgInterval(v pgtype.Interval) string {
	return formatIntervalISO(v)
}
