package validator

import (
	"testing"
)

func TestLoadPackageGoInterfacesErrorAnnotation(t *testing.T) {
	st := &SymbolTable{
		Models: make(map[string]ModelSymbol),
	}
	st.loadPackageGoInterfaces("auth", "../testdata/pkg_error_test")

	// auth._func 모델 키로 등록되어야 함
	ms, ok := st.Models["auth._func"]
	if !ok {
		t.Fatal("auth._func model not found in st.Models")
	}

	// VerifyPassword: @error 401
	mi, ok := ms.Methods["VerifyPassword"]
	if !ok {
		t.Fatal("VerifyPassword method not found")
	}
	if mi.ErrStatus != 401 {
		t.Errorf("VerifyPassword ErrStatus = %d, want 401", mi.ErrStatus)
	}
	// ParamTypes도 수집되어야 함
	if mi.ParamTypes == nil || mi.ParamTypes["Email"] != "string" {
		t.Errorf("VerifyPassword ParamTypes missing Email field")
	}

	// Charge: @error 없음 → 0
	mi2, ok := ms.Methods["Charge"]
	if !ok {
		t.Fatal("Charge method not found")
	}
	if mi2.ErrStatus != 0 {
		t.Errorf("Charge ErrStatus = %d, want 0", mi2.ErrStatus)
	}
}
