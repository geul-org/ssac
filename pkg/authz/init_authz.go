//ff:func feature=pkg-authz type=loader control=sequence
//ff:what 글로벌 인가 상태를 초기화한다 (OPA policy 로딩만; DB 접근 없음)
package authz

import (
	"fmt"
	"os"
)

var globalPolicy string
var globalOwnerships []OwnershipMapping

// Init initializes the global authz state by loading the OPA policy from
// OPA_POLICY_PATH. When DISABLE_AUTHZ=1, Init is a no-op — Check also
// short-circuits to allow.
//
// DB access has been removed: ownership lookups are now the caller's
// responsibility (CheckRequest.Owners). ownerships is retained as a meta-
// table that validators / codegen can consult, but Check itself never
// queries a database.
func Init(policyPath string, ownerships []OwnershipMapping) error {
	globalOwnerships = ownerships

	if os.Getenv("DISABLE_AUTHZ") == "1" {
		return nil
	}

	if policyPath == "" {
		policyPath = os.Getenv("OPA_POLICY_PATH")
	}
	if policyPath == "" {
		return fmt.Errorf("OPA policy path is required (pass argument or set OPA_POLICY_PATH; set DISABLE_AUTHZ=1 to skip)")
	}

	policyData, err := os.ReadFile(policyPath)
	if err != nil {
		return fmt.Errorf("read OPA policy %s: %w", policyPath, err)
	}

	globalPolicy = string(policyData)
	return nil
}

// Ownerships returns the registered @ownership mappings so validators and
// codegen can enumerate required owner-lookup queries without reaching into
// package-private state.
func Ownerships() []OwnershipMapping { return globalOwnerships }
