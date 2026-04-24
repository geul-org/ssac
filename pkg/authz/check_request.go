//ff:type feature=pkg-authz type=model
//ff:what 인가 검사 요청 구조체 — caller 가 Owners 를 사전 로드해 주입
package authz

import "context"

// CheckRequest holds the inputs for an authorization check.
//
// Claim is passed through to OPA input as the `claims` object. Callers should
// use a struct with JSON tags whose keys match the claim names expected by
// the rego policy (e.g. `json:"user_id"`, `json:"org_id"`). The concrete type
// is opaque to this package — rego.Input marshals it via json.Marshal.
//
// Ctx carries the request context used for OPA evaluation. When nil, Check
// falls back to context.Background(). New callers should always propagate the
// handler's ctx.
//
// Owners is the caller-loaded ownership lookup table. Keys are resource
// names (matching the `@ownership <resource>:` declarations in the Rego
// policy) and values are maps from resource-id to owner-id. Both ID types
// are stringified so the policy can compare `data.owners.<resource>[id] ==
// input.claims.user_id` without caring about the underlying column type
// (int64 / string / uuid). The handler is responsible for populating this
// map (typically via a yongol-generated `OwnerLookup<Resource>` sqlc query
// called under the request's pgx.Tx for MVCC-consistent reads).
//
// Since ssac does not touch the database, all DB-dependent fields
// (`*sql.Tx`, `*sql.DB`) have been removed; the caller is the single
// source of authority for DB access.
type CheckRequest struct {
	Ctx        context.Context
	Action     string
	Resource   string
	ResourceID int64
	Claim      any
	Owners     map[string]map[string]string
}
