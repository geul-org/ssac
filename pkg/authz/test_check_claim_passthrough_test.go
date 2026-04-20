//ff:func feature=pkg-authz type=test control=sequence
//ff:what Claim 필드가 OPA input.claims 로 그대로 전달되는지 검증한다
package authz

import (
	"os"
	"testing"
)

// testClaim mirrors a project-local CurrentUser with JSON tags matching rego
// claim keys. OPA json-marshals this struct when building input.
type testClaim struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	OrgID  int64  `json:"org_id"`
}

const passthroughPolicy = `package authz

default allow := false

allow if {
	input.action == "TestOp"
	input.claims.user_id == 1
	input.claims.org_id == 42
	input.claims.role == "admin"
}
`

func TestCheckClaim_StructWithJSONTags(t *testing.T) {
	os.Unsetenv("DISABLE_AUTHZ")
	globalPolicy = passthroughPolicy
	globalDB = nil
	globalOwnerships = nil
	defer func() { globalPolicy = "" }()

	_, err := Check(CheckRequest{
		Action:   "TestOp",
		Resource: "test",
		Claim:    testClaim{UserID: 1, Email: "x@y.z", Role: "admin", OrgID: 42},
	})
	if err != nil {
		t.Fatalf("expected allow, got: %v", err)
	}
}

func TestCheckClaim_PointerStruct(t *testing.T) {
	os.Unsetenv("DISABLE_AUTHZ")
	globalPolicy = passthroughPolicy
	globalDB = nil
	globalOwnerships = nil
	defer func() { globalPolicy = "" }()

	_, err := Check(CheckRequest{
		Action:   "TestOp",
		Resource: "test",
		Claim:    &testClaim{UserID: 1, Email: "x@y.z", Role: "admin", OrgID: 42},
	})
	if err != nil {
		t.Fatalf("pointer passthrough failed: %v", err)
	}
}

func TestCheckClaim_NilDeniesWithClaimsRego(t *testing.T) {
	os.Unsetenv("DISABLE_AUTHZ")
	globalPolicy = passthroughPolicy
	globalDB = nil
	globalOwnerships = nil
	defer func() { globalPolicy = "" }()

	_, err := Check(CheckRequest{
		Action:   "TestOp",
		Resource: "test",
		// Claim: nil (default)
	})
	if err == nil {
		t.Fatal("expected forbidden with nil claim (rego checks user_id/org_id)")
	}
}
