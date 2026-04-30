//ff:func feature=pkg-pgtypex type=util control=sequence
//ff:what ISO 8601 duration 을 pgtype.Interval (Months, Days, Microseconds) 로 파싱한다
package pgtypex

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

// parseIntervalISO accepts the subset PnYnMnDTnHnMnS where each component
// is optional, integer for Y/M/D/H/Min and decimal for S. The result
// matches formatIntervalISO's emission. Errors propagate to the caller —
// the public ToPgInterval wrapper coerces them to Valid:false.
func parseIntervalISO(s string) (pgtype.Interval, error) {
	if !strings.HasPrefix(s, "P") {
		return pgtype.Interval{}, errors.New("pgtypex: ISO 8601 duration must start with 'P'")
	}
	rest := s[1:]
	dateRest := rest
	timeRest := ""
	if i := strings.IndexByte(rest, 'T'); i >= 0 {
		dateRest = rest[:i]
		timeRest = rest[i+1:]
	}
	if dateRest == "" && timeRest == "" {
		return pgtype.Interval{}, errors.New("pgtypex: empty ISO 8601 duration body")
	}
	var months, days int32
	if err := scanIntervalComponents(dateRest, false, &months, &days, nil); err != nil {
		return pgtype.Interval{}, err
	}
	var us int64
	if err := scanIntervalComponents(timeRest, true, nil, nil, &us); err != nil {
		return pgtype.Interval{}, err
	}
	return pgtype.Interval{Months: months, Days: days, Microseconds: us, Valid: true}, nil
}
