✅ 완료

# Phase 7: Generate() 코드젠 잔여 불일치 수정

수정지시서 006 기반. Generate()가 생성하는 서비스 코드의 잔여 불일치 5건을 수정한다.

## 목표

Generate()와 GenerateModelInterfaces()의 시그니처를 완전히 일치시키고, 생성 코드가 Go 컴파일을 통과하도록 한다.

## ssac 범위 vs fullend 범위

| 이슈 | 범위 | 이유 |
|---|---|---|
| 1. FindByID에 opts 잘못 전달 | ssac | buildTemplateData의 opts 추가 조건 |
| 2. ListLessons에 opts 누락 | ssac | 동일 (1과 같은 근본 원인) |
| 3. 리터럴 파라미터 이름 잘못됨 | ssac | resolveLiteralParamName 로직 |
| 4. total이 response에 미포함 | ssac | response json 템플릿 |
| 5. sqlc 쿼리 이름 접두사 | ssac | loadSqlcQueries 파싱 |
| 6. 패키지/디렉토리 불일치 | **fullend** | outDir 결정은 fullend 책임 |
| 7. model. 타입 한정자 | **fullend** | server.go 생성은 fullend 책임 |

## 작업 순서

### Step 1: opts 전달을 List 메서드로 한정 (이슈 1, 2 통합)

**문제**: 현재 opts 추가는 서비스 함수의 `HasQueryOpts()`로 판단 → 같은 함수 내 모든 get 시퀀스에 opts가 추가됨. FindByID에도 opts가 들어가고, 다른 함수의 List 메서드에는 opts가 빠짐.

**수정**: opts 전달을 **모델 메서드 이름** 기준으로 판단.

`generator/go_target.go` `buildTemplateData()` 수정:

```go
// 현재: 함수 레벨 HasQueryOpts 체크
if st != nil && seq.Type == parser.SeqGet {
    if op, ok := st.Operations[funcName]; ok && op.HasQueryOpts() {
        d.ParamArgs += ", opts"
        d.HasTotal = true  // if []result
    }
}

// 변경: 모델 메서드명 기준 체크
if st != nil && seq.Type == parser.SeqGet && seq.Model != "" {
    if isListMethod(seq.Model) && hasAnyQueryOpts(st) {
        if d.ParamArgs != "" {
            d.ParamArgs += ", "
        }
        d.ParamArgs += "opts"
        if seq.Result != nil && strings.HasPrefix(seq.Result.Type, "[]") {
            d.HasTotal = true
        }
    }
}
```

`isListMethod()` 신규:
```go
func isListMethod(model string) bool {
    parts := strings.SplitN(model, ".", 2)
    if len(parts) < 2 { return false }
    return strings.HasPrefix(parts[1], "List")
}
```

`hasAnyQueryOpts()` 신규: 어떤 operation이든 QueryOpts가 있으면 true (List 메서드는 항상 opts를 받을 수 있도록).

`GenerateFunc()` 내 `opts := QueryOpts{}` 생성도 동일 조건으로 변경:
```go
// 현재: 함수 레벨 HasQueryOpts
// 변경: 시퀀스 중 List 메서드가 있고 QueryOpts가 존재하면 생성
needsOpts := false
for _, seq := range sf.Sequences {
    if seq.Type == parser.SeqGet && isListMethod(seq.Model) && hasAnyQueryOpts(st) {
        needsOpts = true
        break
    }
}
if needsOpts {
    buf.WriteString("\topts := QueryOpts{}\n\n")
}
```

### Step 2: DDLTable에 컬럼 순서 추가 (이슈 3 전제)

**문제**: `DDLTable.Columns`는 `map[string]string`으로 순서가 없음. 리터럴 파라미터의 DDL 역매핑에서 남은 컬럼을 알파벳 순으로 선택하면 잘못된 컬럼이 선택됨 (password_hash < role, method < status).

**수정**: `DDLTable`에 `ColumnOrder []string` 필드 추가.

`validator/symbol.go`:
```go
type DDLTable struct {
    Columns     map[string]string
    ColumnOrder []string    // DDL 정의 순서 보존
    ForeignKeys []ForeignKey
    Indexes     []Index
}
```

`parseDDLTables()`에서 컬럼 파싱 시 순서 보존:
```go
t.Columns[colName] = goType
t.ColumnOrder = append(t.ColumnOrder, colName)
```

### Step 3: 리터럴 파라미터 이름 결정 개선 (이슈 3)

**문제**: `resolveLiteralParamName()`이 남은 string 컬럼 중 알파벳 순 첫 번째를 선택 → `password_hash`가 `role`보다 먼저 선택됨.

**수정 1**: usedColumns 구축 개선 — dot notation에서 복합 컬럼명도 추가:
```go
// 현재: enrollment.ID → usedColumns["id"]
// 변경: enrollment.ID → usedColumns["id"] + usedColumns["enrollment_id"]
if strings.Contains(p.Name, ".") {
    parts := strings.SplitN(p.Name, ".", 2)
    usedColumns[toSnakeCase(parts[1])] = true
    usedColumns[toSnakeCase(parts[0])+"_"+toSnakeCase(parts[1])] = true
}
```

**수정 2**: `resolveLiteralParamName()`을 알파벳순 → DDL 컬럼 순서로 변경:
```go
func resolveLiteralParamName(modelName string, usedColumns map[string]bool, st *validator.SymbolTable) string {
    tableName := toSnakeCase(modelName) + "s"
    table, ok := st.DDLTables[tableName]
    if !ok { return "" }

    autoColumns := map[string]bool{
        "id": true, "created_at": true, "updated_at": true, "deleted_at": true,
    }

    // DDL 정의 순서로 순회하여 첫 번째 미사용 string 컬럼 선택
    for _, col := range table.ColumnOrder {
        goType := table.Columns[col]
        if autoColumns[col] || usedColumns[col] || goType != "string" {
            continue
        }
        return lcFirst(snakeToCamel(col))
    }
    return ""
}
```

이렇게 하면 DDL 순서상 `role`이 `password_hash`보다 뒤에 있어도... 음, 실제 DDL 순서에 따라 다를 수 있다. 핵심은 알파벳순이 아닌 DDL 정의 순서를 따르는 것이다.

**추가**: password_hash가 이미 쓰인 것으로 판단되려면, `hashedPassword` 변수 참조가 `password_hash` 컬럼과 매칭되어야 한다. 현재 `toSnakeCase("hashedPassword")` = `hashed_password` ≠ `password_hash`. 이를 위해 usedColumns 구축 시 DDL 역조회를 추가:

```go
// 변수 참조 (@param hashedPassword): DDL 컬럼과 직접 매칭 시도
// toSnakeCase(name)이 DDL에 없으면, 해당 모델 테이블에서 유사 컬럼 탐색
if p.Source == "" && !strings.Contains(p.Name, ".") && !strings.HasPrefix(p.Name, `"`) {
    snake := toSnakeCase(p.Name)
    usedColumns[snake] = true
    // DDL 직접 조회: 해당 모델 테이블에서 snake가 없으면 부분 매칭
    tableName := toSnakeCase(modelName) + "s"
    if table, ok := st.DDLTables[tableName]; ok {
        if _, exists := table.Columns[snake]; !exists {
            // 부분 매칭: "hashed_password" → "password_hash"는 매칭 안 됨
            // 대안: name의 주요 단어가 포함된 컬럼 제외
            // "hashedPassword" → words: "hashed", "password"
            // DDL에 "password" 포함 컬럼 → "password_hash" → 제외
            words := splitCamelWords(p.Name)
            for col := range table.Columns {
                for _, w := range words {
                    if strings.Contains(col, strings.ToLower(w)) {
                        usedColumns[col] = true
                    }
                }
            }
        }
    }
}
```

`splitCamelWords()` 신규: `"hashedPassword"` → `["hashed", "Password"]`

### Step 4: response json에 total 자동 포함 (이슈 4)

**문제**: List 3-tuple에서 `total` 변수를 받지만 response JSON에 포함하지 않아 Go 컴파일 에러 (미사용 변수).

**수정**: `templateData`에 `HasTotal` 정보를 response까지 전달. GenerateFunc에서 함수 내 HasTotal 여부를 추적하고, response json 템플릿에서 조건부 추가.

`generator/go_target.go` — `templateData`에 이미 `HasTotal`이 있으므로, response 시퀀스에도 전달:

```go
// GenerateFunc 내: 함수 레벨 hasTotal 추적
funcHasTotal := false
for i, seq := range sf.Sequences {
    data := buildTemplateData(seq, &errDeclared, resultTypes, st, sf.Name)
    if data.HasTotal {
        funcHasTotal = true
    }
    // response 시퀀스에 funcHasTotal 전달
    if strings.HasPrefix(seq.Type, "response") {
        data.HasTotal = funcHasTotal
    }
    // ... 템플릿 실행
}
```

`generator/go_templates.go` — response json 템플릿 수정:

```
{{- define "response json" -}}
	// response json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		{{- range .Vars}}
		"{{.}}": {{.}},
		{{- end}}
		{{- if .HasTotal}}
		"total": total,
		{{- end}}
	})
{{end}}
```

### Step 5: sqlc 쿼리 이름 접두사 지원 (이슈 5)

**문제**: sqlc에서 `FindByID`가 여러 파일에 중복되면 충돌. 해결책으로 `CourseFindByID` 형태의 접두사를 사용하지만, ssac은 `FindByID`만 인식.

**수정**: `loadSqlcQueries()`에서 쿼리 이름을 파싱할 때, 모델명 접두사를 분리.

`validator/symbol.go` `loadSqlcQueries()` 수정:

```go
for scanner.Scan() {
    line := strings.TrimSpace(scanner.Text())
    if strings.HasPrefix(line, "-- name:") {
        parts := strings.Fields(line)
        if len(parts) >= 4 {
            queryName := parts[2]
            cardinality := strings.TrimPrefix(parts[3], ":")
            // 모델명 접두사 분리: "CourseFindByID" → "FindByID"
            methodName := stripModelPrefix(queryName, modelName)
            ms.Methods[methodName] = MethodInfo{Cardinality: cardinality}
        }
    }
}
```

`stripModelPrefix()` 신규:
```go
func stripModelPrefix(queryName, modelName string) string {
    if strings.HasPrefix(queryName, modelName) {
        stripped := queryName[len(modelName):]
        if len(stripped) > 0 && stripped[0] >= 'A' && stripped[0] <= 'Z' {
            return stripped
        }
    }
    return queryName  // 접두사 없으면 원본 유지
}
```

이렇게 하면 `CourseFindByID` → `FindByID`, `CourseList` → `List`로 분리. 접두사가 없는 기존 쿼리도 그대로 동작.

### Step 6: 테스트 업데이트

기존 테스트 유지 + 변경/추가:

1. `TestGenerateQueryOptsAndTotal` — total이 response JSON에 포함되는지 확인 추가
2. `TestGenerateModelInterfaces` — opts가 List 메서드에만 추가되는지 확인
3. `TestStripModelPrefix` — sqlc 접두사 분리 단위 테스트 (validator_test.go)
4. DDL 컬럼 순서 보존 확인 (validator_test.go)

## 변경 파일 목록

| 파일 | 변경 유형 |
|---|---|
| `generator/go_target.go` | 수정: opts를 List 메서드 한정, resolveLiteralParamName DDL 순서 + 부분 매칭, response에 HasTotal 전달, isListMethod/hasAnyQueryOpts/splitCamelWords 신규 |
| `generator/go_templates.go` | 수정: response json 템플릿에 total 조건부 추가 |
| `generator/generator_test.go` | 수정: QueryOpts/total 테스트 업데이트 |
| `validator/symbol.go` | 수정: DDLTable에 ColumnOrder 추가, parseDDLTables 순서 보존, loadSqlcQueries 접두사 분리, stripModelPrefix 신규 |
| `validator/validator_test.go` | 수정: DDL 컬럼 순서 + sqlc 접두사 테스트 추가 |

## 하지 않는 것

- 이슈 6 (패키지/디렉토리 불일치): fullend 측에서 outDir을 조정하여 해결
- 이슈 7 (model. 타입 한정자): fullend의 server.go 생성 로직에서 해결
- QueryOpts 실제 바인딩 (`r.URL.Query()` → QueryOpts 필드): 후속 Phase

## 검증

```bash
go test ./parser/... ./validator/... ./generator/... -count=1

# dummy-study gen 재생성으로 확인:
# 1. get_reservation.go: FindByID(reservationID) — opts 없음
# 2. list_my_reservations.go: ListByUserID(currentUser.UserID, opts) — opts 있음
# 3. list_my_reservations.go: response에 "total": total 포함
# 4. models_gen.go: Session.Create(userID int64) — dot notation 유지
```

## 리스크

- **중간**: DDLTable에 ColumnOrder 추가는 validator 구조체 변경이지만, 기존 필드는 유지하므로 하위 호환
- 리터럴 파라미터의 부분 매칭(splitCamelWords)이 과도하게 제외할 가능성 있음 → 테스트로 검증
- sqlc 접두사 분리는 기존 접두사 없는 쿼리와 호환 (stripModelPrefix가 매칭 안 되면 원본 유지)
