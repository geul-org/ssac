package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/ssac/artifacts/internal/parser"
)

// Generate는 []ServiceFunc를 받아 outDir에 Go 파일을 생성한다.
func Generate(funcs []parser.ServiceFunc, outDir string) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("출력 디렉토리 생성 실패: %w", err)
	}

	for _, sf := range funcs {
		code, err := GenerateFunc(sf)
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
func GenerateFunc(sf parser.ServiceFunc) ([]byte, error) {
	var buf bytes.Buffer
	imports := collectImports(sf.Sequences)

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
	fmt.Fprintf(&buf, "func %s(w http.ResponseWriter, r *http.Request) {\n", sf.Name)

	// request 파라미터 추출
	requestParams := collectRequestParams(sf.Sequences)
	for _, p := range requestParams {
		fmt.Fprintf(&buf, "\t%s := r.FormValue(%q)\n", lcFirst(p), p)
	}
	if len(requestParams) > 0 {
		buf.WriteString("\n")
	}

	// sequence 블록 생성
	errDeclared := false
	for i, seq := range sf.Sequences {
		data := buildTemplateData(seq, &errDeclared)

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
	Target string
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

func buildTemplateData(seq parser.Sequence, errDeclared *bool) templateData {
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

	// guard 대상
	d.Target = seq.Target

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

// collectImports는 사용된 패키지를 수집한다.
func collectImports(seqs []parser.Sequence) []string {
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

	var imports []string
	// 표준 라이브러리 먼저
	order := []string{"encoding/json", "net/http", "golang.org/x/crypto/bcrypt"}
	for _, imp := range order {
		if seen[imp] {
			imports = append(imports, imp)
		}
	}
	return imports
}

// collectRequestParams는 source가 "request"인 고유 파라미터명을 수집한다.
func collectRequestParams(seqs []parser.Sequence) []string {
	seen := map[string]bool{}
	var params []string
	for _, seq := range seqs {
		for _, p := range seq.Params {
			if p.Source == "request" && !seen[p.Name] {
				seen[p.Name] = true
				params = append(params, p.Name)
			}
		}
	}
	return params
}

// buildParamArgs는 Param 슬라이스를 함수 호출 인자 문자열로 변환한다.
func buildParamArgs(params []parser.Param) string {
	var args []string
	for _, p := range params {
		args = append(args, resolveParamRef(p.Name))
	}
	return strings.Join(args, ", ")
}

// resolveParamRef는 파라미터 참조를 Go 표현식으로 변환한다.
// "ProjectID" (request) → lcFirst → "projectID" (이미 추출된 변수)
// "project.OwnerEmail" → 그대로
// "\"리터럴\"" → 그대로
func resolveParamRef(name string) string {
	if name == "" {
		return ""
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

// lcFirst는 첫 글자를 소문자로 변환한다.
func lcFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}
