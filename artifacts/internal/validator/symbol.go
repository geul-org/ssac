package validator

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// SymbolTable은 외부 SSOT에서 수집한 심볼 정보다.
type SymbolTable struct {
	Models     map[string]ModelSymbol     // "User" → {Methods: {"FindByID": true}}
	Operations map[string]OperationSymbol // "Login" → {RequestFields, ResponseFields}
	Components map[string]bool            // "notification" → true
	Funcs      map[string]bool            // "calculateRefund" → true
}

// ModelSymbol은 모델의 메서드 목록이다.
type ModelSymbol struct {
	Methods map[string]bool
}

// OperationSymbol은 API 엔드포인트의 request/response 필드 목록이다.
type OperationSymbol struct {
	RequestFields  map[string]bool
	ResponseFields map[string]bool
}

// LoadSymbolTable은 프로젝트 디렉토리에서 심볼 테이블을 구성한다.
// 디렉토리 구조:
//
//	<root>/db/queries/*.sql  — sqlc 쿼리 (모델+메서드)
//	<root>/api/openapi.yaml  — OpenAPI spec (request/response)
//	<root>/model/*.go        — Go interface (component, func)
func LoadSymbolTable(root string) (*SymbolTable, error) {
	st := &SymbolTable{
		Models:     make(map[string]ModelSymbol),
		Operations: make(map[string]OperationSymbol),
		Components: make(map[string]bool),
		Funcs:      make(map[string]bool),
	}

	if err := st.loadSqlcQueries(filepath.Join(root, "db", "queries")); err != nil {
		return nil, fmt.Errorf("sqlc 쿼리 로드 실패: %w", err)
	}
	if err := st.loadOpenAPI(filepath.Join(root, "api", "openapi.yaml")); err != nil {
		return nil, fmt.Errorf("OpenAPI 로드 실패: %w", err)
	}
	if err := st.loadGoInterfaces(filepath.Join(root, "model")); err != nil {
		return nil, fmt.Errorf("Go interface 로드 실패: %w", err)
	}

	return st, nil
}

// loadSqlcQueries는 queries/*.sql에서 모델과 메서드를 추출한다.
// 파일명: users.sql → 모델 "User" (단수화 + PascalCase)
// 주석: -- name: FindByID :one → 메서드 "FindByID"
func (st *SymbolTable) loadSqlcQueries(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		modelName := sqlFileToModel(entry.Name())
		ms := ModelSymbol{Methods: make(map[string]bool)}

		f, err := os.Open(filepath.Join(dir, entry.Name()))
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			// -- name: FindByID :one
			if strings.HasPrefix(line, "-- name:") {
				parts := strings.Fields(line)
				if len(parts) >= 3 {
					ms.Methods[parts[2]] = true
				}
			}
		}
		f.Close()

		if len(ms.Methods) > 0 {
			st.Models[modelName] = ms
		}
	}
	return nil
}

// sqlFileToModel은 "reservations.sql" → "Reservation" 변환한다.
func sqlFileToModel(filename string) string {
	name := strings.TrimSuffix(filename, ".sql")
	// 단수화: 간단한 규칙 (es → 제거, s → 제거)
	if strings.HasSuffix(name, "ies") {
		name = name[:len(name)-3] + "y"
	} else if strings.HasSuffix(name, "ses") || strings.HasSuffix(name, "xes") {
		name = name[:len(name)-2]
	} else if strings.HasSuffix(name, "s") {
		name = name[:len(name)-1]
	}
	// PascalCase
	return strings.ToUpper(name[:1]) + name[1:]
}

// loadOpenAPI는 openapi.yaml에서 operationId별 request/response 필드를 추출한다.
func (st *SymbolTable) loadOpenAPI(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var spec openAPISpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return fmt.Errorf("YAML 파싱 실패: %w", err)
	}

	schemas := spec.Components.Schemas

	for _, pathItem := range spec.Paths {
		for _, op := range pathItem.operations() {
			if op.OperationID == "" {
				continue
			}

			opSym := OperationSymbol{
				RequestFields:  make(map[string]bool),
				ResponseFields: make(map[string]bool),
			}

			// path/query parameters
			for _, param := range op.Parameters {
				opSym.RequestFields[param.Name] = true
			}

			// request body fields
			if op.RequestBody != nil {
				if content, ok := op.RequestBody.Content["application/json"]; ok {
					fields := collectSchemaFields(content.Schema, schemas)
					for _, f := range fields {
						opSym.RequestFields[f] = true
					}
				}
			}

			// response fields (200)
			if resp, ok := op.Responses["200"]; ok {
				if content, ok := resp.Content["application/json"]; ok {
					fields := collectSchemaFields(content.Schema, schemas)
					for _, f := range fields {
						opSym.ResponseFields[f] = true
					}
				}
			}

			st.Operations[op.OperationID] = opSym
		}
	}

	return nil
}

// loadGoInterfaces는 model/*.go에서 interface(component)와 func을 추출한다.
func (st *SymbolTable) loadGoInterfaces(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	fset := token.NewFileSet()
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		f, err := parser.ParseFile(fset, filepath.Join(dir, entry.Name()), nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("%s 파싱 실패: %w", entry.Name(), err)
		}

		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				// interface → component로 등록 (소문자 이름)
				if _, ok := ts.Type.(*ast.InterfaceType); ok {
					componentName := strings.ToLower(ts.Name.Name[:1]) + ts.Name.Name[1:]
					st.Components[componentName] = true

					// interface의 메서드도 Models에 등록
					ms := ModelSymbol{Methods: make(map[string]bool)}
					iface := ts.Type.(*ast.InterfaceType)
					for _, method := range iface.Methods.List {
						if len(method.Names) > 0 {
							ms.Methods[method.Names[0].Name] = true
						}
					}
					if len(ms.Methods) > 0 {
						st.Models[ts.Name.Name] = ms
					}
				}
			}

			// 패키지 레벨 func → Funcs로 등록
			fd, ok := decl.(*ast.FuncDecl)
			if ok && fd.Recv == nil {
				st.Funcs[fd.Name.Name] = true
			}
		}

		// ast.Decls에서 FuncDecl은 GenDecl과 별개로 순회
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if ok && fd.Recv == nil {
				st.Funcs[fd.Name.Name] = true
			}
		}
	}

	return nil
}

// --- OpenAPI YAML 구조체 ---

type openAPISpec struct {
	Paths      map[string]openAPIPathItem `yaml:"paths"`
	Components openAPIComponents          `yaml:"components"`
}

type openAPIComponents struct {
	Schemas map[string]openAPISchema `yaml:"schemas"`
}

type openAPISchema struct {
	Type       string                   `yaml:"type"`
	Properties map[string]openAPISchema `yaml:"properties"`
	Ref        string                   `yaml:"$ref"`
}

type openAPIPathItem struct {
	Get    *openAPIOperation `yaml:"get"`
	Post   *openAPIOperation `yaml:"post"`
	Put    *openAPIOperation `yaml:"put"`
	Delete *openAPIOperation `yaml:"delete"`
}

func (p openAPIPathItem) operations() []*openAPIOperation {
	var ops []*openAPIOperation
	for _, op := range []*openAPIOperation{p.Get, p.Post, p.Put, p.Delete} {
		if op != nil {
			ops = append(ops, op)
		}
	}
	return ops
}

type openAPIOperation struct {
	OperationID string                       `yaml:"operationId"`
	Parameters  []openAPIParameter           `yaml:"parameters"`
	RequestBody *openAPIRequestBody          `yaml:"requestBody"`
	Responses   map[string]openAPIResponse   `yaml:"responses"`
}

type openAPIParameter struct {
	Name string `yaml:"name"`
	In   string `yaml:"in"`
}

type openAPIRequestBody struct {
	Content map[string]openAPIMediaType `yaml:"content"`
}

type openAPIResponse struct {
	Content map[string]openAPIMediaType `yaml:"content"`
}

type openAPIMediaType struct {
	Schema openAPISchema `yaml:"schema"`
}

// collectSchemaFields는 인라인 properties와 $ref 모두에서 필드를 수집한다.
func collectSchemaFields(schema openAPISchema, schemas map[string]openAPISchema) []string {
	var fields []string

	// 인라인 properties
	for k := range schema.Properties {
		fields = append(fields, k)
	}

	// $ref 해결
	if schema.Ref != "" {
		name := schema.Ref[strings.LastIndex(schema.Ref, "/")+1:]
		if resolved, ok := schemas[name]; ok {
			for k := range resolved.Properties {
				fields = append(fields, k)
			}
		}
	}

	return fields
}
