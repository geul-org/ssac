// Package pgtypex provides deterministic bridge functions between
// OpenAPI-surfaced Go types (oapi-codegen `runtime/types`, stdlib `string`,
// `int64`, `time.Time`, …) and their sqlc / pgx v5 pgtype wrappers
// (`pgtype.UUID`, `pgtype.Text`, `pgtype.Int8`, `pgtype.Timestamptz`, …).
//
// Each supported PG type exposes five helpers:
//
//	ToPg<T>(v) <Pg>          // NOT NULL wrap (Valid: true)
//	ToPg<T>Ptr(p) <Pg>       // nullable wrap (nil → Valid: false)
//	FromPg<T>(v) <Api>       // NOT NULL unwrap; caller has guaranteed Valid
//	FromPg<T>Ptr(v) *<Api>   // nullable unwrap (!Valid → nil)
//	IsNilPg<T>(v) bool       // !v.Valid alias for nil-check sites
//
// And one bulk helper:
//
//	ToPg<T>s([]<Api>) []<Pg> // slice conversion for sqlc bulk params
//
// UUID additionally exposes `UUIDToString` for OPA Owners-map sites where
// the UUID must be serialised as a canonical 8-4-4-4-12 hex string.
//
// # Design notes
//
// The package contains no business decisions. Function bodies are
// 1-3 lines of mechanical wrap / unwrap. The reason this lives in a
// runtime library — rather than emitted into each generated SaaS — is
// that the conversion logic is schema-independent: pgx v5 and oapi-codegen
// fix the layout of every type involved, so there is nothing project-specific
// to preserve. Centralising here removes per-SaaS boilerplate, lets a
// pgx / oapi-codegen major upgrade be handled once, and keeps generated
// `internal/service/` directories free of preserve-hash-tracked scaffolding.
//
// # Version pin discipline
//
// pgx v5 and oapi-codegen runtime types are imported transitively through
// this package. SaaS go.mod files MUST pin the same major versions ssac
// pins, otherwise `pgtype.UUID` from a different major would be a different
// Go type and the bridges would not type-check at the call site. yongol
// validate enforces this with a "bridge import drift" rule.
//
// # Invariants
//
// FromPg<T> functions assume the caller has already established v.Valid.
// In yongol-generated code this guarantee comes from the DDL `NOT NULL`
// constraint plus validate rule D-12. Calling FromPg<T> on an invalid
// pgtype value returns the zero value of <Api> — silent rather than
// panicking, because a panic would break the response cycle of an
// upstream HTTP handler that already passed validation.
package pgtypex
