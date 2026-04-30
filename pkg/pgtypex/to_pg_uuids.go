//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what []OpenAPI UUID 를 []pgtype.UUID 로 bulk wrap 한다 (nil → nil)
package pgtypex

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oapi-codegen/runtime/types"
)

func ToPgUUIDs(vs []types.UUID) []pgtype.UUID {
	if vs == nil {
		return nil
	}
	out := make([]pgtype.UUID, len(vs))
	for i, v := range vs {
		out[i] = pgtype.UUID{Bytes: v, Valid: true}
	}
	return out
}
