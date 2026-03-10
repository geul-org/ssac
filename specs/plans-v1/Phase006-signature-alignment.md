✅ 완료

# Phase 6: Generate()와 GenerateModelInterfaces() 시그니처 정렬

수정지시서 005 기반. 서비스 함수 코드와 모델 인터페이스 간 시그니처 불일치를 수정한다.

## 목표

`Generate()`가 생성하는 메서드 호출과 `GenerateModelInterfaces()`가 생성하는 인터페이스의 시그니처를 일치시켜 Go 컴파일이 통과하도록 한다.

## 문제 3건

### 문제 1: List 계열 메서드에 QueryOpts 미전달

`GenerateModelInterfaces()`는 x-확장 있으면 `opts QueryOpts`를 인터페이스에 추가하지만, `Generate()`는 해당 인자를 전달하지 않음.

### 문제 2: 파생 파라미터/리터럴이 인터페이스에 누락

`resolveParamName()`이 dot notation(`enrollment.ID`)과 리터럴(`"pending"`)을 빈 문자열로 반환 → `renderParams()`에서 스킵됨.

### 문제 3: List 계열 반환값 불일치

인터페이스는 `many`+QueryOpts면 `([]T, int, error)` 3-tuple이지만, `Generate()`는 항상 `result, err :=` 2-tuple 생성.

## 작업 순서

### Step 1: 파생 파라미터/리터럴 인터페이스 포함 (문제 2)

**변경 파일**: `generator/go_target.go`

`resolveParamName()` 수정:

| 현재 | 변경 후 |
|---|---|
| dot notation (`enrollment.ID`) → `""` (스킵) | → `enrollmentID` (부분 결합, lcFirst) |
| 리터럴 (`"pending"`) → `""` (스킵) | → DDL 컬럼 역매핑으로 이름 결정 |

dot notation 처리:
```go
// "enrollment.ID" → "enrollmentID"
// "course.Price" → "coursePrice"
func resolveParamName(p parser.Param) string {
    if strings.HasPrefix(p.Name, `"`) {
        // 리터럴 → 별도 처리
        return ""  // 이 시점에서는 빈 이름, 후처리로 해결
    }
    if strings.Contains(p.Name, ".") {
        parts := strings.SplitN(p.Name, ".", 2)
        return parts[0] + strings.Title(parts[1])  // enrollment + ID
        // 단, strings.Title deprecated → 수동 첫글자 대문자
    }
    return lcFirst(p.Name)
}
```

`resolveParamType()` 수정 — dot notation 지원:
```go
// "enrollment.ID" → Enrollment 테이블의 id 컬럼 타입 조회
// "course.Price" → Course 테이블의 price 컬럼 타입 조회
func resolveParamType(p parser.Param, modelName string, st *validator.SymbolTable) string {
    if strings.HasPrefix(p.Name, `"`) {
        return "string"
    }
    if strings.Contains(p.Name, ".") {
        parts := strings.SplitN(p.Name, ".", 2)
        refTable := toSnakeCase(parts[0]) + "s"
        refCol := toSnakeCase(parts[1])
        if table, ok := st.DDLTables[refTable]; ok {
            if goType, ok := table.Columns[refCol]; ok {
                return goType
            }
        }
        return "string"  // fallback
    }
    // ... 기존 로직
}
```

리터럴 파라미터 처리 — DDL 컬럼 역매핑:
- 해당 모델의 DDL 테이블 컬럼 목록에서, 이미 다른 @param으로 사용된 컬럼을 제외하고 남는 컬럼 매칭 시도
- 매칭 불가 시 positional 이름 (`arg0`, `arg1`, ...) 사용
- 구현: `deriveInterfaces()` 내에서 리터럴 파라미터에 대해 DDL 역매핑 수행

```go
// 접근: 해당 메서드의 모든 param 중 이름 있는 것의 snake_case를 수집
// DDL 테이블 컬럼에서 이 이름들을 제외
// 남은 string 타입 컬럼 중 하나를 리터럴에 매핑
```

### Step 2: QueryOpts 전달 (문제 1)

**변경 파일**: `generator/go_target.go`

현재 `buildTemplateData()`는 `seq`, `errDeclared`, `resultTypes`만 받아 st/funcName에 접근 불가.

방안: `buildTemplateData()`에 st와 funcName 파라미터 추가.

```go
func buildTemplateData(seq parser.Sequence, errDeclared *bool, resultTypes map[string]string, st *validator.SymbolTable, funcName string) templateData {
    // ... 기존 로직 ...
    d.ParamArgs = buildParamArgs(seq.Params)

    // QueryOpts 추가: get 시퀀스 + HasQueryOpts
    if st != nil && (seq.Type == parser.SeqGet) {
        if op, ok := st.Operations[funcName]; ok && op.HasQueryOpts() {
            if d.ParamArgs != "" {
                d.ParamArgs += ", "
            }
            d.ParamArgs += "opts"
        }
    }
}
```

QueryOpts 구성 코드 생성 — `GenerateFunc()` 내에서 함수 본문 상단에 추가:

```go
// GenerateFunc 내부, request 파라미터 추출 후, sequence 블록 전에:
if st != nil {
    if op, ok := st.Operations[sf.Name]; ok && op.HasQueryOpts() {
        buf.WriteString(generateQueryOptsCode(op))
    }
}
```

`generateQueryOptsCode()` 신규 함수:
```go
func generateQueryOptsCode(op validator.OperationSymbol) string {
    // 빈 QueryOpts{} 전달 (단순화 옵션 채택)
    // 향후 실제 query parameter 바인딩은 별도 Phase
    return "\topts := QueryOpts{}\n\n"
}
```

**단순화 결정**: 수정지시서 권장대로 우선 빈 `QueryOpts{}`를 전달한다. 실제 `r.URL.Query()` 바인딩은 후속 작업.

### Step 3: List 반환값 3-tuple (문제 3)

**변경 파일**: `generator/go_target.go`, `generator/go_templates.go`

`templateData`에 필드 추가:
```go
type templateData struct {
    // ... 기존 ...
    HasTotal bool   // many + QueryOpts → 3-tuple 반환
}
```

`buildTemplateData()` 수정:
```go
// get 시퀀스 + many cardinality + HasQueryOpts → HasTotal = true
if st != nil && seq.Type == parser.SeqGet && seq.Result != nil {
    if strings.HasPrefix(seq.Result.Type, "[]") {
        if op, ok := st.Operations[funcName]; ok && op.HasQueryOpts() {
            d.HasTotal = true
        }
    }
}
```

`go_templates.go` get 템플릿 수정:
```
{{- define "get" -}}
	// get
	{{if .HasTotal}}{{.Result.Var}}, total, err{{else}}{{.Result.Var}}, err{{end}} := {{.ModelVar}}.{{.ModelMethod}}({{.ParamArgs}})
	if err != nil {
		http.Error(w, "{{.Message}}", http.StatusInternalServerError)
		return
	}
{{end}}
```

response 템플릿에서 total 포함 — `templateData`에 `HasTotal` 있으면 response에 total 추가:
- response json 템플릿에서 total을 별도로 추가할지는 보류 (response는 @var로 명시하므로 spec에서 `@var total` 추가하면 됨)
- 이 Phase에서는 get 시퀀스의 3-tuple 반환만 해결

### Step 4: 호출부 시그니처 업데이트

`GenerateFunc()` 내에서 `buildTemplateData()` 호출 시그니처 변경:
```go
// 기존
data := buildTemplateData(seq, &errDeclared, resultTypes)
// 변경
data := buildTemplateData(seq, &errDeclared, resultTypes, st, sf.Name)
```

### Step 5: 테스트 추가

기존 테스트 유지 + 새 테스트 케이스:

1. **QueryOpts 전달 테스트**: dummy-study의 `list_reservations.go` 생성 → `opts` 인자 포함 확인
2. **3-tuple 반환 테스트**: `reservations, total, err :=` 패턴 확인
3. **파생 파라미터 인터페이스 테스트**: `TestGenerateModelInterfaces`에 dot notation/리터럴 포함 체크 추가
   - 현재 `"no dot notation", "Create() (*Token, error)"` 체크 → Session.Create는 dot notation만 사용 (user.ID) → 인터페이스에 `userID`로 포함되어야 함
   - 이 체크를 `"Session Create with dot param", "Create(userID int64) (*Token, error)"` 등으로 변경

**주의**: 기존 테스트 `"no dot notation", "Create() (*Token, error)"` 은 Session.Create의 dot notation param(user.ID)이 스킵되는 현재 동작을 테스트한다. 이 수정 후에는 `user.ID` → `userID int64`로 인터페이스에 포함되므로 체크를 변경해야 한다.

## 변경 파일 목록

| 파일 | 변경 유형 |
|---|---|
| `generator/go_target.go` | 수정: resolveParamName/resolveParamType dot notation 지원, buildTemplateData에 st/funcName 추가, QueryOpts/HasTotal 처리, generateQueryOptsCode 신규 |
| `generator/go_templates.go` | 수정: get 템플릿에 HasTotal 조건부 3-tuple |
| `generator/generator_test.go` | 수정: 기존 체크 업데이트 + QueryOpts/3-tuple/파생 파라미터 테스트 추가 |

## 검증

```bash
go test ./generator/... -count=1

# dummy-study 재생성으로 시그니처 일치 확인
# 1. list_reservations.go: reservationModel.ListByUserID(currentUser.UserID, opts) 호출
# 2. list_reservations.go: reservations, total, err := (3-tuple)
# 3. models_gen.go: Session.Create(userID int64) — dot notation 포함
# 4. models_gen.go: 리터럴 파라미터 컬럼명 포함
```

## 리스크

- **중간**: 기능 변경이므로 기존 테스트 수정 필요 (Session.Create 인터페이스 시그니처 변경)
- 리터럴 파라미터 DDL 역매핑이 실패할 수 있음 → positional fallback으로 대응
- QueryOpts 구성은 빈 struct로 단순화 → 컴파일은 통과하나 런타임 동작은 미완성
