//ff:func feature=pkg-authz type=test control=sequence
//ff:what CheckRequest 구조체 필드가 올바르게 설정되는지 검증한다
package authz

import "testing"

func TestCheckRequestFields(t *testing.T) {
	claim := map[string]any{"user_id": int64(42), "role": "client"}
	req := CheckRequest{
		Action:     "AcceptProposal",
		Resource:   "gig",
		ResourceID: "99",
		Claim:      claim,
	}
	if req.Action != "AcceptProposal" {
		t.Fatal("Action mismatch")
	}
	if req.Resource != "gig" {
		t.Fatal("Resource mismatch")
	}
	if req.ResourceID != "99" {
		t.Fatal("ResourceID mismatch")
	}
	got, ok := req.Claim.(map[string]any)
	if !ok {
		t.Fatal("Claim type assertion failed")
	}
	if got["user_id"] != int64(42) {
		t.Fatalf("user_id mismatch: %v", got["user_id"])
	}
	if got["role"] != "client" {
		t.Fatalf("role mismatch: %v", got["role"])
	}
}
