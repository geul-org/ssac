package generator

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/park-jun-woo/ssac/artifacts/internal/parser"
)

func specsDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "specs", "backend", "service")
}

func TestGenerateCreateSession(t *testing.T) {
	sf, err := parser.ParseFile(filepath.Join(specsDir(), "create_session.go"))
	if err != nil {
		t.Fatalf("파싱 실패: %v", err)
	}

	code, err := GenerateFunc(*sf)
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

	code, err := GenerateFunc(*sf)
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
		{"guard exists", "if sessionCount != nil"},
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

	code, err := GenerateFunc(*sf)
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
	if err := Generate(funcs, outDir); err != nil {
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

func readFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	return string(b), err
}
