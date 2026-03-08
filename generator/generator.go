package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"

	"github.com/geul-org/ssac/parser"
	"github.com/geul-org/ssac/validator"
)

// Generate는 []ServiceFunc를 받아 outDir에 Go 파일을 생성한다.
// st가 non-nil이면 DDL 타입 기반 변환 코드를 생성한다.
func Generate(funcs []parser.ServiceFunc, outDir string, st *validator.SymbolTable) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("출력 디렉토리 생성 실패: %w", err)
	}

	for _, sf := range funcs {
		code, err := GenerateFunc(sf, st)
		if err != nil {
			return fmt.Errorf("%s 코드 생성 실패: %w", sf.Name, err)
		}

		path := filepath.Join(outDir, sf.FileName)
		if err := os.WriteFile(path, code, 0644); err != nil {
			return fmt.Errorf("%s 파일 쓰기 실패: %w", path, err)
		}
	}
	return nil
}

// GenerateFunc는 단일 ServiceFunc의 Go 코드를 생성한다.
// st가 non-nil이면 DDL 타입 기반 변환 코드를 생성한다.
func GenerateFunc(sf parser.ServiceFunc, st *validator.SymbolTable) ([]byte, error) {
	var buf bytes.Buffer

	// path parameter 결정
	var pathParams []validator.PathParam
	if st != nil {
		if op, ok := st.Operations[sf.Name]; ok {
			pathParams = op.PathParams
		}
	}
	pathParamSet := map[string]bool{}
	for _, pp := range pathParams {
		pathParamSet[pp.Name] = true
	}

	// request 파라미터 타입 결정 (path param은 제외)
	typedParams := collectTypedRequestParams(sf.Sequences, st, pathParamSet)
	imports := collectImports(sf.Sequences, typedParams)

	// package
	buf.WriteString("package service\n\n")

	// imports
	if len(imports) > 0 {
		buf.WriteString("import (\n")
		for _, imp := range imports {
			fmt.Fprintf(&buf, "\t%q\n", imp)
		}
		buf.WriteString(")\n\n")
	}

	// func signature
	sig := "func %s(w http.ResponseWriter, r *http.Request"
	if len(pathParams) > 0 {
		var ppArgs []string
		for _, pp := range pathParams {
			ppArgs = append(ppArgs, fmt.Sprintf("%s %s", lcFirst(pp.Name), pp.GoType))
		}
		sig += ", " + strings.Join(ppArgs, ", ")
	}
	sig += ") {\n"
	fmt.Fprintf(&buf, sig, sf.Name)

	// request 파라미터 추출 (타입 변환 포함, path param 제외)
	for _, tp := range typedParams {
		buf.WriteString(tp.extractCode)
	}
	if len(typedParams) > 0 {
		buf.WriteString("\n")
	}

	// result 타입 맵 구축 (guard 비교식 결정용)
	resultTypes := map[string]string{}
	for _, seq := range sf.Sequences {
		if seq.Result != nil {
			resultTypes[seq.Result.Var] = seq.Result.Type
		}
	}

	// sequence 블록 생성
	// 타입 변환 코드에서 err를 선언했으면 이미 선언된 것으로 처리
	errDeclared := hasConversionErr(typedParams)
	for i, seq := range sf.Sequences {
		data := buildTemplateData(seq, &errDeclared, resultTypes)

		tmplName := templateName(seq)
		var seqBuf bytes.Buffer
		if err := templates.ExecuteTemplate(&seqBuf, tmplName, data); err != nil {
			return nil, fmt.Errorf("sequence[%d] %s 템플릿 실행 실패: %w", i, seq.Type, err)
		}
		buf.Write(seqBuf.Bytes())
		buf.WriteString("\n")
	}

	buf.WriteString("}\n")

	// gofmt
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.Bytes(), fmt.Errorf("gofmt 실패: %w\n--- raw ---\n%s", err, buf.String())
	}
	return formatted, nil
}

// templateData는 템플릿에 전달하는 데이터다.
type templateData struct {
	// 공통
	Message string
	// get, post, put, delete
	ModelVar    string
	ModelMethod string
	ParamArgs   string
	Result      *parser.Result
	// guard
	Target      string
	ZeroCheck   string // "== nil", "== 0", `== ""`, "== false"
	ExistsCheck string // "!= nil", "> 0", `!= ""`, "== true"
	// authorize
	Action   string
	Resource string
	ID       string
	// password
	Hash  string
	Plain string
	// call
	Component       string
	ComponentMethod string
	Func            string
	FirstErr        bool
	// response
	Vars []string
}

func buildTemplateData(seq parser.Sequence, errDeclared *bool, resultTypes map[string]string) templateData {
	d := templateData{
		Message: seq.Message,
		Result:  seq.Result,
		Vars:    seq.Vars,
	}

	// 모델 분리: "Project.FindByID" → "projectModel", "FindByID"
	if seq.Model != "" {
		parts := strings.SplitN(seq.Model, ".", 2)
		d.ModelVar = lcFirst(parts[0]) + "Model"
		if len(parts) > 1 {
			d.ModelMethod = parts[1]
		}
	}

	// 기본 메시지 생성
	if d.Message == "" {
		d.Message = defaultMessage(seq)
	}

	// 파라미터 인자 문자열
	d.ParamArgs = buildParamArgs(seq.Params)

	// guard 대상 + 타입별 비교식
	d.Target = seq.Target
	if seq.Type == parser.SeqGuardNil || seq.Type == parser.SeqGuardExists {
		typeName := resultTypes[seq.Target]
		d.ZeroCheck, d.ExistsCheck = zeroValueChecks(typeName)
	}

	// authorize
	d.Action = seq.Action
	d.Resource = seq.Resource
	d.ID = resolveParamRef(seq.ID)

	// password
	if seq.Type == parser.SeqPassword && len(seq.Params) >= 2 {
		d.Hash = resolveParamRef(seq.Params[0].Name)
		d.Plain = resolveParamRef(seq.Params[1].Name)
	}

	// call
	d.Component = seq.Component
	d.ComponentMethod = "Execute"
	d.Func = seq.Func

	// err 선언 추적
	// 좌변에 새 변수가 있으면 := (Go 규칙: 하나라도 새 변수면 := 가능)
	// authorize: allowed 새 변수 → 항상 :=
	// get, post: result 새 변수 → 항상 :=
	// put, delete: result 없음 → err만 → 첫 사용 시 :=, 이후 =
	// call: result 있으면 := (새 변수), 없으면 err만 → 첫 사용 시 :=, 이후 =
	switch seq.Type {
	case parser.SeqAuthorize, parser.SeqGet, parser.SeqPost:
		d.FirstErr = true // 항상 새 변수가 좌변에 있음
		*errDeclared = true
	case parser.SeqCall:
		if seq.Result != nil {
			d.FirstErr = true // 새 변수 있음
			*errDeclared = true
		} else if !*errDeclared {
			d.FirstErr = true
			*errDeclared = true
		}
	case parser.SeqPut, parser.SeqDelete:
		if !*errDeclared {
			d.FirstErr = true
			*errDeclared = true
		}
	}

	return d
}

func templateName(seq parser.Sequence) string {
	if seq.Type == parser.SeqCall {
		if seq.Component != "" {
			return "call_component"
		}
		return "call_func"
	}
	// response json → "response json"
	if strings.HasPrefix(seq.Type, "response") {
		return seq.Type
	}
	return seq.Type
}

// typedRequestParam은 request 파라미터의 타입과 추출 코드를 보관한다.
type typedRequestParam struct {
	name        string // PascalCase 원본명
	goType      string // "string", "int64", "time.Time" 등
	extractCode string // 추출 코드 (줄바꿈 포함)
}

// collectTypedRequestParams는 source가 "request"인 파라미터를 수집하고 DDL 타입을 결정한다.
// pathParamSet에 포함된 파라미터는 함수 인자로 이미 받으므로 제외한다.
// request 파라미터가 2개 이상이면 JSON body로 간주하여 struct decode 코드를 생성한다.
func collectTypedRequestParams(seqs []parser.Sequence, st *validator.SymbolTable, pathParamSet map[string]bool) []typedRequestParam {
	seen := map[string]bool{}
	var rawParams []struct {
		name   string
		goType string
	}
	for _, seq := range seqs {
		for _, p := range seq.Params {
			if p.Source != "request" || seen[p.Name] || pathParamSet[p.Name] {
				continue
			}
			seen[p.Name] = true

			goType := "string"
			if st != nil {
				goType = lookupDDLType(p.Name, st)
			}
			rawParams = append(rawParams, struct {
				name   string
				goType string
			}{p.Name, goType})
		}
	}

	// 심볼 테이블이 있고 request 파라미터가 2개 이상이면 JSON body struct decode
	if st != nil && len(rawParams) >= 2 {
		return buildJSONBodyParams(rawParams)
	}

	// 1개 이하면 기존 FormValue 방식
	var params []typedRequestParam
	for _, rp := range rawParams {
		varName := lcFirst(rp.name)
		code := generateExtractCode(varName, rp.name, rp.goType)
		params = append(params, typedRequestParam{
			name:        rp.name,
			goType:      rp.goType,
			extractCode: code,
		})
	}
	return params
}

// buildJSONBodyParams는 JSON body struct decode + 변수 추출 코드를 생성한다.
// 여러 타입이 사용될 수 있으므로 별도의 typedRequestParam 엔트리로 import 힌트를 추가한다.
func buildJSONBodyParams(rawParams []struct {
	name   string
	goType string
}) []typedRequestParam {
	var buf bytes.Buffer

	// struct 정의
	buf.WriteString("\tvar req struct {\n")
	for _, rp := range rawParams {
		jsonTag := toSnakeCase(rp.name)
		buf.WriteString(fmt.Sprintf("\t\t%s %s `json:\"%s\"`\n", rp.name, rp.goType, jsonTag))
	}
	buf.WriteString("\t}\n")

	// decode
	buf.WriteString("\tif err := json.NewDecoder(r.Body).Decode(&req); err != nil {\n")
	buf.WriteString("\t\thttp.Error(w, \"invalid request body\", http.StatusBadRequest)\n")
	buf.WriteString("\t\treturn\n")
	buf.WriteString("\t}\n")

	// 변수 추출
	for _, rp := range rawParams {
		varName := lcFirst(rp.name)
		buf.WriteString(fmt.Sprintf("\t%s := req.%s\n", varName, rp.name))
	}

	// 전체 코드를 json_body 엔트리에 담고, time.Time import 힌트도 추가
	result := []typedRequestParam{{
		name:        "_json_body",
		goType:      "json_body",
		extractCode: buf.String(),
	}}
	// time.Time은 struct 필드에 사용되므로 import 필요 (strconv 등은 불필요)
	for _, rp := range rawParams {
		if rp.goType == "time.Time" {
			result = append(result, typedRequestParam{
				name:   rp.name,
				goType: rp.goType,
			})
			break // 한 번만 추가하면 충분
		}
	}
	return result
}

// lookupDDLType은 PascalCase 파라미터명을 snake_case로 변환하여 DDL 컬럼 타입을 조회한다.
func lookupDDLType(paramName string, st *validator.SymbolTable) string {
	snakeName := toSnakeCase(paramName)
	for _, table := range st.DDLTables {
		if goType, ok := table.Columns[snakeName]; ok {
			return goType
		}
	}
	return "string"
}

// generateExtractCode는 타입별 request 파라미터 추출 코드를 생성한다.
func generateExtractCode(varName, paramName, goType string) string {
	switch goType {
	case "int64":
		return fmt.Sprintf("\t%s, err := strconv.ParseInt(r.FormValue(%q), 10, 64)\n"+
			"\tif err != nil {\n"+
			"\t\thttp.Error(w, \"%s: 유효하지 않은 값\", http.StatusBadRequest)\n"+
			"\t\treturn\n"+
			"\t}\n", varName, paramName, paramName)
	case "float64":
		return fmt.Sprintf("\t%s, err := strconv.ParseFloat(r.FormValue(%q), 64)\n"+
			"\tif err != nil {\n"+
			"\t\thttp.Error(w, \"%s: 유효하지 않은 값\", http.StatusBadRequest)\n"+
			"\t\treturn\n"+
			"\t}\n", varName, paramName, paramName)
	case "bool":
		return fmt.Sprintf("\t%s, err := strconv.ParseBool(r.FormValue(%q))\n"+
			"\tif err != nil {\n"+
			"\t\thttp.Error(w, \"%s: 유효하지 않은 값\", http.StatusBadRequest)\n"+
			"\t\treturn\n"+
			"\t}\n", varName, paramName, paramName)
	case "time.Time":
		return fmt.Sprintf("\t%s, err := time.Parse(time.RFC3339, r.FormValue(%q))\n"+
			"\tif err != nil {\n"+
			"\t\thttp.Error(w, \"%s: 유효하지 않은 값\", http.StatusBadRequest)\n"+
			"\t\treturn\n"+
			"\t}\n", varName, paramName, paramName)
	default: // string
		return fmt.Sprintf("\t%s := r.FormValue(%q)\n", varName, paramName)
	}
}

// collectImports는 사용된 패키지를 수집한다.
func collectImports(seqs []parser.Sequence, typedParams []typedRequestParam) []string {
	seen := map[string]bool{
		"net/http": true, // 항상 사용
	}

	for _, seq := range seqs {
		switch {
		case strings.HasPrefix(seq.Type, "response json"):
			seen["encoding/json"] = true
		case seq.Type == parser.SeqPassword:
			seen["golang.org/x/crypto/bcrypt"] = true
		}
	}

	for _, tp := range typedParams {
		switch tp.goType {
		case "int64", "float64", "bool":
			seen["strconv"] = true
		case "time.Time":
			seen["time"] = true
		case "json_body":
			seen["encoding/json"] = true
		}
	}

	var imports []string
	// 표준 라이브러리 먼저, 알파벳 순
	order := []string{"encoding/json", "net/http", "strconv", "time", "golang.org/x/crypto/bcrypt"}
	for _, imp := range order {
		if seen[imp] {
			imports = append(imports, imp)
		}
	}
	return imports
}

// buildParamArgs는 Param 슬라이스를 함수 호출 인자 문자열로 변환한다.
func buildParamArgs(params []parser.Param) string {
	var args []string
	for _, p := range params {
		args = append(args, resolveParam(p))
	}
	return strings.Join(args, ", ")
}

// resolveParam은 Param의 source를 고려하여 Go 표현식으로 변환한다.
func resolveParam(p parser.Param) string {
	// 예약어 source → source.Name
	if p.Source == "currentUser" || p.Source == "config" {
		return p.Source + "." + p.Name
	}
	return resolveParamRef(p.Name)
}

// resolveParamRef는 파라미터 참조를 Go 표현식으로 변환한다.
// "ProjectID" (request) → lcFirst → "projectID" (이미 추출된 변수)
// "project.OwnerEmail" → 그대로
// "\"리터럴\"" → 그대로
// "new" → nil (리소스 생성 전 authorize에서 ID가 없음을 의미)
func resolveParamRef(name string) string {
	if name == "" {
		return ""
	}
	// "new"는 아직 존재하지 않는 리소스를 의미 → nil
	if name == "new" {
		return "nil"
	}
	// 따옴표 리터럴은 그대로
	if strings.HasPrefix(name, `"`) {
		return name
	}
	// dot notation은 그대로
	if strings.Contains(name, ".") {
		return name
	}
	return lcFirst(name)
}

// extractGuardTarget은 guard 시퀀스에서 대상 변수를 추출한다.
// "@sequence guard nil project" → Type에서 3번째 단어가 아니므로,
// 바로 앞 sequence의 result를 사용하거나, Params[0]을 사용한다.
// 현재 구현: message에서 유추하거나, 직전 get의 result를 쓰는 대신
// spec 파일에서 guard nil 뒤에 오는 단어를 Target으로 쓴다.
// 이를 위해 parser에서 guard의 대상을 추출해야 하지만, 현재 Type에만 포함되어 있다.
// → parser를 수정하지 않고, Sequence에 Target 필드를 추가하는 대신
//
//	여기서는 Type 파싱 시 잘린 부분을 활용한다.
//
// 실제로는 parser가 guard target을 별도로 저장해야 하므로, 여기서 처리한다.
// defaultMessage는 sequence 타입과 모델명으로 기본 에러 메시지를 생성한다.
func defaultMessage(seq parser.Sequence) string {
	modelName := ""
	if seq.Model != "" {
		parts := strings.SplitN(seq.Model, ".", 2)
		modelName = parts[0]
	}

	switch seq.Type {
	case parser.SeqGet:
		return modelName + " 조회 실패"
	case parser.SeqPost:
		return modelName + " 생성 실패"
	case parser.SeqPut:
		return modelName + " 수정 실패"
	case parser.SeqDelete:
		return modelName + " 삭제 실패"
	case parser.SeqGuardNil:
		return seq.Target + "가 존재하지 않습니다"
	case parser.SeqGuardExists:
		return seq.Target + "가 이미 존재합니다"
	case parser.SeqAuthorize:
		return "권한 확인 실패"
	case parser.SeqPassword:
		return "비밀번호가 일치하지 않습니다"
	case parser.SeqCall:
		if seq.Component != "" {
			return seq.Component + " 호출 실패"
		}
		if seq.Func != "" {
			return seq.Func + " 호출 실패"
		}
		return "호출 실패"
	}
	return "처리 실패"
}

// zeroValueChecks는 타입에 따른 guard 비교식을 반환한다.
func zeroValueChecks(typeName string) (zeroCheck, existsCheck string) {
	switch typeName {
	case "int", "int32", "int64", "float64":
		return "== 0", "> 0"
	case "bool":
		return "== false", "== true"
	case "string":
		return `== ""`, `!= ""`
	default:
		return "== nil", "!= nil"
	}
}

// hasConversionErr는 타입 변환이 있어서 err가 이미 선언되었는지 반환한다.
func hasConversionErr(params []typedRequestParam) bool {
	for _, p := range params {
		// json_body는 if err := 로 스코프 내 선언이므로 제외
		if p.goType != "string" && p.goType != "json_body" {
			return true
		}
	}
	return false
}

// lcFirst는 첫 글자를 소문자로 변환한다.
func lcFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}
