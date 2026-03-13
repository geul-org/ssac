package generator

import (
	"bytes"
	"fmt"
	"go/format"

	"github.com/geul-org/ssac/parser"
	"github.com/geul-org/ssac/validator"
)

// generateHTTPFunc는 HTTP 핸들러 함수를 생성한다.
func (g *GoTarget) generateHTTPFunc(sf parser.ServiceFunc, st *validator.SymbolTable) ([]byte, error) {
	var buf bytes.Buffer

	// 분석
	pathParams := getPathParams(sf.Name, st)
	pathParamSet := map[string]bool{}
	for _, pp := range pathParams {
		pathParamSet[pp.Name] = true
	}

	requestParams := collectRequestParams(sf.Sequences, st, pathParamSet)
	needsCU := needsCurrentUser(sf.Sequences)
	needsQO := needsQueryOpts(sf, st)
	imports := collectImports(sf, requestParams, pathParams, needsCU, needsQO)

	// package
	pkgName := "service"
	if sf.Domain != "" {
		pkgName = sf.Domain
	}
	buf.WriteString("package " + pkgName + "\n\n")

	// imports
	if len(imports) > 0 {
		buf.WriteString("import (\n")
		for _, imp := range imports {
			fmt.Fprintf(&buf, "\t%q\n", imp)
		}
		buf.WriteString(")\n\n")
	}

	// func signature
	fmt.Fprintf(&buf, "func (h *Handler) %s(c *gin.Context) {\n", sf.Name)

	// path parameters
	for _, pp := range pathParams {
		buf.WriteString(generatePathParamCode(pp))
	}
	if len(pathParams) > 0 {
		buf.WriteString("\n")
	}

	// currentUser
	if needsCU {
		var cuBuf bytes.Buffer
		goTemplates.ExecuteTemplate(&cuBuf, "currentUser", nil)
		buf.Write(cuBuf.Bytes())
		buf.WriteString("\n")
	}

	// request parameters
	for _, rp := range requestParams {
		buf.WriteString(rp.extractCode)
	}
	if len(requestParams) > 0 {
		buf.WriteString("\n")
	}

	// QueryOpts
	if needsQO {
		buf.WriteString(generateQueryOptsCode(sf.Name, st))
		buf.WriteString("\n")
	}

	// Transaction
	useTx := hasWriteSequence(sf.Sequences)
	if useTx {
		buf.WriteString("\ttx, err := h.DB.BeginTx(c.Request.Context(), nil)\n")
		buf.WriteString("\tif err != nil {\n")
		buf.WriteString("\t\tc.JSON(http.StatusInternalServerError, gin.H{\"error\": \"transaction failed\"})\n")
		buf.WriteString("\t\treturn\n")
		buf.WriteString("\t}\n")
		buf.WriteString("\tdefer tx.Rollback()\n\n")
	}

	// result types for guard checks
	resultTypes := map[string]string{}
	for _, seq := range sf.Sequences {
		if seq.Result != nil {
			resultTypes[seq.Result.Var] = seq.Result.Type
		}
	}

	// sequences
	errDeclared := hasConversionErr(requestParams)
	if useTx {
		errDeclared = true
	}
	declaredVars := map[string]bool{}
	funcHasTotal := false
	usedVars := collectUsedVars(sf.Sequences)
	committed := false
	for i, seq := range sf.Sequences {
		if useTx && seq.Type == parser.SeqResponse && !committed {
			buf.WriteString("\tif err = tx.Commit(); err != nil {\n")
			buf.WriteString("\t\tc.JSON(http.StatusInternalServerError, gin.H{\"error\": \"commit failed\"})\n")
			buf.WriteString("\t\treturn\n")
			buf.WriteString("\t}\n\n")
			committed = true
		}
		data := buildTemplateData(seq, &errDeclared, declaredVars, resultTypes, st, sf.Name, useTx)
		if data.HasTotal {
			funcHasTotal = true
		}
		if seq.Type == parser.SeqResponse {
			data.HasTotal = funcHasTotal
		}
		// 미사용 변수 처리
		if seq.Result != nil && !usedVars[seq.Result.Var] {
			data.Unused = true
			if data.ErrDeclared {
				data.ReAssign = true // _, err = (no new vars with :=)
			}
		}

		tmplName := templateName(seq)
		var seqBuf bytes.Buffer
		if err := goTemplates.ExecuteTemplate(&seqBuf, tmplName, data); err != nil {
			return nil, fmt.Errorf("sequence[%d] %s 템플릿 실행 실패: %w", i, seq.Type, err)
		}
		buf.Write(seqBuf.Bytes())
		buf.WriteString("\n")
	}

	if useTx && !committed {
		buf.WriteString("\tif err = tx.Commit(); err != nil {\n")
		buf.WriteString("\t\tc.JSON(http.StatusInternalServerError, gin.H{\"error\": \"commit failed\"})\n")
		buf.WriteString("\t\treturn\n")
		buf.WriteString("\t}\n\n")
	}

	buf.WriteString("}\n")

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.Bytes(), fmt.Errorf("gofmt 실패: %w\n--- raw ---\n%s", err, buf.String())
	}
	return formatted, nil
}

// generateSubscribeFunc는 큐 구독 핸들러 함수를 생성한다.
func (g *GoTarget) generateSubscribeFunc(sf parser.ServiceFunc, st *validator.SymbolTable) ([]byte, error) {
	var buf bytes.Buffer

	pkgName := "service"
	if sf.Domain != "" {
		pkgName = sf.Domain
	}
	buf.WriteString("package " + pkgName + "\n\n")

	imports := collectSubscribeImports(sf)
	if len(imports) > 0 {
		buf.WriteString("import (\n")
		for _, imp := range imports {
			fmt.Fprintf(&buf, "\t%q\n", imp)
		}
		buf.WriteString(")\n\n")
	}

	msgType := sf.Subscribe.MessageType
	fmt.Fprintf(&buf, "func (h *Handler) %s(ctx context.Context, message %s) error {\n", sf.Name, msgType)

	resultTypes := map[string]string{}
	for _, seq := range sf.Sequences {
		if seq.Result != nil {
			resultTypes[seq.Result.Var] = seq.Result.Type
		}
	}

	useTx := hasWriteSequence(sf.Sequences)
	if useTx {
		buf.WriteString("\ttx, err := h.DB.BeginTx(ctx, nil)\n")
		buf.WriteString("\tif err != nil {\n")
		buf.WriteString("\t\treturn fmt.Errorf(\"transaction failed: %w\", err)\n")
		buf.WriteString("\t}\n")
		buf.WriteString("\tdefer tx.Rollback()\n\n")
	}

	errDeclared := useTx
	declaredVars := map[string]bool{}
	usedVars := collectUsedVars(sf.Sequences)
	for i, seq := range sf.Sequences {
		data := buildTemplateData(seq, &errDeclared, declaredVars, resultTypes, st, sf.Name, useTx)
		// 미사용 변수 처리
		if seq.Result != nil && !usedVars[seq.Result.Var] {
			data.Unused = true
			if data.ErrDeclared {
				data.ReAssign = true // _, err = (no new vars with :=)
			}
		}
		tmplName := subscribeTemplateName(seq)
		var seqBuf bytes.Buffer
		if err := goTemplates.ExecuteTemplate(&seqBuf, tmplName, data); err != nil {
			return nil, fmt.Errorf("sequence[%d] %s 템플릿 실행 실패: %w", i, seq.Type, err)
		}
		buf.Write(seqBuf.Bytes())
		buf.WriteString("\n")
	}

	if useTx {
		buf.WriteString("\tif err = tx.Commit(); err != nil {\n")
		buf.WriteString("\t\treturn fmt.Errorf(\"commit failed: %w\", err)\n")
		buf.WriteString("\t}\n\n")
	}

	buf.WriteString("\treturn nil\n")
	buf.WriteString("}\n")

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.Bytes(), fmt.Errorf("gofmt 실패: %w\n--- raw ---\n%s", err, buf.String())
	}
	return formatted, nil
}

// subscribeTemplateName은 subscribe 함수 내 시퀀스의 템플릿 이름을 반환한다.
func subscribeTemplateName(seq parser.Sequence) string {
	switch seq.Type {
	case parser.SeqCall:
		if seq.Result != nil {
			return "sub_call_with_result"
		}
		return "sub_call_no_result"
	case parser.SeqPublish:
		return "sub_publish"
	default:
		return "sub_" + seq.Type
	}
}

// templateName은 HTTP 함수 내 시퀀스의 템플릿 이름을 반환한다.
func templateName(seq parser.Sequence) string {
	switch seq.Type {
	case parser.SeqResponse:
		if seq.Target != "" {
			return "response_direct"
		}
		return "response"
	case parser.SeqCall:
		if seq.Result != nil {
			return "call_with_result"
		}
		return "call_no_result"
	case parser.SeqPublish:
		return "publish"
	default:
		return seq.Type
	}
}
