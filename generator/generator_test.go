package generator

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/geul-org/ssac/parser"
	"github.com/geul-org/ssac/validator"
)

func specsDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "testdata", "backend-service")
}

func TestGenerateCreateSession(t *testing.T) {
	sf, err := parser.ParseFile(filepath.Join(specsDir(), "create_session.go"))
	if err != nil {
		t.Fatalf("파싱 실패: %v", err)
	}

	code, err := GenerateFunc(*sf, nil)
	if err != nil {
		t.Fatalf("코드 생성 실패: %v", err)
	}

	got := string(code)

	// 핵심 코드 블록 존재 확인
	checks := []struct {
		label string
		want  string
	}{
		{"package", "package service"},
		{"import json", `"encoding/json"`},
		{"import http", `"net/http"`},
		{"func sig", "func CreateSession(w http.ResponseWriter, r *http.Request)"},
		{"request param", `projectID := r.FormValue("ProjectID")`},
		{"request param", `command := r.FormValue("Command")`},
		{"get model call", "projectModel.FindByID(projectID)"},
		{"get result", "project, err :="},
		{"guard nil", "if project == nil"},
		{"guard nil message", `"프로젝트가 존재하지 않습니다"`},
		{"post model call", "sessionModel.Create(projectID, command)"},
		{"post result", "session, err :="},
		{"response json", "json.NewEncoder(w).Encode"},
		{"response var", `"session": session`},
	}

	for _, c := range checks {
		if !strings.Contains(got, c.want) {
			t.Errorf("[%s] %q 없음\n--- got ---\n%s", c.label, c.want, got)
		}
	}
}

func TestGenerateDeleteProject(t *testing.T) {
	sf, err := parser.ParseFile(filepath.Join(specsDir(), "delete_project.go"))
	if err != nil {
		t.Fatalf("파싱 실패: %v", err)
	}

	code, err := GenerateFunc(*sf, nil)
	if err != nil {
		t.Fatalf("코드 생성 실패: %v", err)
	}

	got := string(code)

	checks := []struct {
		label string
		want  string
	}{
		{"authorize check", `authz.Check(currentUser, "delete", "project", projectID)`},
		{"authorize err", "allowed, err :="},
		{"authorize forbidden", "if !allowed"},
		{"get project", "projectModel.FindByID(projectID)"},
		{"guard nil project", "if project == nil"},
		{"get sessionCount", "sessionModel.CountByProjectID(projectID)"},
		{"guard exists", "if sessionCount > 0"},
		{"guard exists msg", "하위 세션이 존재하여 삭제할 수 없습니다"},
		{"call component", "notification"},
		{"call func", "cleanupProjectResources(project)"},
		{"delete", "projectModel.Delete(projectID)"},
	}

	for _, c := range checks {
		if !strings.Contains(got, c.want) {
			t.Errorf("[%s] %q 없음\n--- got ---\n%s", c.label, c.want, got)
		}
	}
}

func TestGenerateGofmt(t *testing.T) {
	sf, err := parser.ParseFile(filepath.Join(specsDir(), "create_session.go"))
	if err != nil {
		t.Fatalf("파싱 실패: %v", err)
	}

	code, err := GenerateFunc(*sf, nil)
	if err != nil {
		t.Fatalf("코드 생성 실패: %v", err)
	}

	// gofmt가 적용되면 탭 인덴트가 있어야 함
	if !strings.Contains(string(code), "\t") {
		t.Error("gofmt 적용 안 됨: 탭 인덴트 없음")
	}
}

func TestGenerateDir(t *testing.T) {
	funcs, err := parser.ParseDir(specsDir())
	if err != nil {
		t.Fatalf("파싱 실패: %v", err)
	}

	outDir := t.TempDir()
	if err := Generate(funcs, outDir, nil); err != nil {
		t.Fatalf("Generate 실패: %v", err)
	}

	// 파일 생성 확인
	for _, sf := range funcs {
		path := filepath.Join(outDir, sf.FileName)
		if _, err := readFile(path); err != nil {
			t.Errorf("파일 생성 안 됨: %s", path)
		}
	}
}

func TestGenerateModelInterfaces(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	dummyRoot := filepath.Join(filepath.Dir(file), "..", "specs", "dummy-study")

	funcs, err := parser.ParseDir(filepath.Join(dummyRoot, "service"))
	if err != nil {
		t.Fatalf("파싱 실패: %v", err)
	}

	st, err := validator.LoadSymbolTable(dummyRoot)
	if err != nil {
		t.Fatalf("심볼 테이블 로드 실패: %v", err)
	}

	outDir := filepath.Join(filepath.Dir(file), "..", "testdata", "model_iface_test")
	os.MkdirAll(outDir, 0755)
	defer os.RemoveAll(outDir)

	if err := GenerateModelInterfaces(funcs, st, outDir); err != nil {
		t.Fatalf("모델 인터페이스 생성 실패: %v", err)
	}

	code, err := os.ReadFile(filepath.Join(outDir, "model", "models_gen.go"))
	if err != nil {
		t.Fatalf("생성된 파일 읽기 실패: %v", err)
	}
	got := string(code)

	checks := []struct {
		label string
		want  string
	}{
		{"package", "package model"},
		{"time import", `import "time"`},
		{"ReservationModel", "type ReservationModel interface"},
		{"Create with time", "startAt time.Time, endAt time.Time"},
		{"FindByID", "FindByID(reservationID int64) (*Reservation, error)"},
		{"ListByUserID many+pagination", "ListByUserID(userID int64, opts QueryOpts) ([]Reservation, int, error)"},
		{"UpdateStatus exec", "UpdateStatus(reservationID int64) error"},
		{"RoomModel", "type RoomModel interface"},
		{"Room Update", "Update(roomID int64, name string, capacity int64, location string) error"},
		{"UserModel", "type UserModel interface"},
		{"FindByEmail", "FindByEmail(email string) (*User, error)"},
		{"SessionModel", "type SessionModel interface"},
		// dot notation params (user.ID) should NOT appear in interface
		{"no dot notation", "Create() (*Token, error)"},
	}

	// SSaC에서 사용되지 않는 메서드는 포함되면 안 됨
	negChecks := []struct {
		label  string
		reject string
	}{
		{"unused ListAll", "ListAll("},
	}

	for _, c := range checks {
		if !strings.Contains(got, c.want) {
			t.Errorf("[%s] %q 없음\n--- got ---\n%s", c.label, c.want, got)
		}
	}
	for _, c := range negChecks {
		if strings.Contains(got, c.reject) {
			t.Errorf("[%s] %q가 포함되면 안 됨", c.label, c.reject)
		}
	}
}

func TestGenerateTypedRequestParams(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	dummyRoot := filepath.Join(filepath.Dir(file), "..", "specs", "dummy-study")

	st, err := validator.LoadSymbolTable(dummyRoot)
	if err != nil {
		t.Fatalf("심볼 테이블 로드 실패: %v", err)
	}

	// create_reservation.go: StartAt, EndAt은 time.Time
	sf, err := parser.ParseFile(filepath.Join(dummyRoot, "service", "create_reservation.go"))
	if err != nil {
		t.Fatalf("파싱 실패: %v", err)
	}

	code, err := GenerateFunc(*sf, st)
	if err != nil {
		t.Fatalf("코드 생성 실패: %v", err)
	}
	got := string(code)

	// create_reservation: 3개 request param (RoomID, StartAt, EndAt) → JSON body
	checks := []struct {
		label string
		want  string
	}{
		{"json import", `"encoding/json"`},
		{"time import", `"time"`},
		{"json body struct", "var req struct"},
		{"StartAt field", "`json:\"start_at\"`"},
		{"EndAt field", "`json:\"end_at\"`"},
		{"json decode", "json.NewDecoder(r.Body).Decode(&req)"},
		{"startAt var", "startAt := req.StartAt"},
		{"roomID var", "roomID := req.RoomID"},
	}
	for _, c := range checks {
		if !strings.Contains(got, c.want) {
			t.Errorf("[%s] %q 없음\n--- got ---\n%s", c.label, c.want, got)
		}
	}

	// FormValue는 사용되지 않아야 함 (JSON body 모드)
	if strings.Contains(got, "r.FormValue") {
		t.Errorf("JSON body 모드에서 r.FormValue가 사용되면 안 됨\n--- got ---\n%s", got)
	}

	// update_room.go: RoomID는 path param, Name/Capacity/Location은 JSON body
	sf2, err := parser.ParseFile(filepath.Join(dummyRoot, "service", "update_room.go"))
	if err != nil {
		t.Fatalf("파싱 실패: %v", err)
	}

	code2, err := GenerateFunc(*sf2, st)
	if err != nil {
		t.Fatalf("코드 생성 실패: %v", err)
	}
	got2 := string(code2)

	checks2 := []struct {
		label string
		want  string
	}{
		{"json import", `"encoding/json"`},
		{"RoomID path param", "func UpdateRoom(w http.ResponseWriter, r *http.Request, roomID int64)"},
		{"json body struct", "var req struct"},
		{"Capacity field", "Capacity int64"},
		{"json decode", "json.NewDecoder(r.Body).Decode(&req)"},
		{"capacity var", "capacity := req.Capacity"},
	}
	for _, c := range checks2 {
		if !strings.Contains(got2, c.want) {
			t.Errorf("[%s] %q 없음\n--- got ---\n%s", c.label, c.want, got2)
		}
	}

	// RoomID는 path param이므로 FormValue/JSON body에 없어야 함
	if strings.Contains(got2, `r.FormValue("RoomID")`) {
		t.Errorf("RoomID는 path param이므로 FormValue가 없어야 함\n--- got ---\n%s", got2)
	}
}

func TestGeneratePathParamSignature(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	dummyRoot := filepath.Join(filepath.Dir(file), "..", "specs", "dummy-study")

	st, err := validator.LoadSymbolTable(dummyRoot)
	if err != nil {
		t.Fatalf("심볼 테이블 로드 실패: %v", err)
	}

	// get_reservation.go: ReservationID는 OpenAPI path param (int64)
	sf, err := parser.ParseFile(filepath.Join(dummyRoot, "service", "get_reservation.go"))
	if err != nil {
		t.Fatalf("파싱 실패: %v", err)
	}

	code, err := GenerateFunc(*sf, st)
	if err != nil {
		t.Fatalf("코드 생성 실패: %v", err)
	}
	got := string(code)

	checks := []struct {
		label string
		want  string
	}{
		{"path param sig", "func GetReservation(w http.ResponseWriter, r *http.Request, reservationID int64)"},
	}

	negChecks := []struct {
		label  string
		reject string
	}{
		{"no FormValue for path param", `r.FormValue("ReservationID")`},
	}

	for _, c := range checks {
		if !strings.Contains(got, c.want) {
			t.Errorf("[%s] %q 없음\n--- got ---\n%s", c.label, c.want, got)
		}
	}
	for _, c := range negChecks {
		if strings.Contains(got, c.reject) {
			t.Errorf("[%s] %q가 포함되면 안 됨\n--- got ---\n%s", c.label, c.reject, got)
		}
	}

	// login.go: path param 없으므로 기존 시그니처 유지
	sf2, err := parser.ParseFile(filepath.Join(dummyRoot, "service", "login.go"))
	if err != nil {
		t.Fatalf("파싱 실패: %v", err)
	}

	code2, err := GenerateFunc(*sf2, st)
	if err != nil {
		t.Fatalf("코드 생성 실패: %v", err)
	}
	got2 := string(code2)

	if !strings.Contains(got2, "func Login(w http.ResponseWriter, r *http.Request)") {
		t.Errorf("Login 시그니처에 path param이 없어야 함\n--- got ---\n%s", got2)
	}
}

func readFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	return string(b), err
}
