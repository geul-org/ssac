# Phase 016: 패키지 접두사 @model의 Go interface 교차 검증

## 목표

`@get Session session = session.Session.Get({token: request.Token})` 처럼 패키지 접두사가 있는 모델 호출을 파싱하고, 해당 패키지의 Go interface와 교차 검증한다.

- 접두사 없음 (`Gig.FindByID`) → 기존 DDL 기반 DB 모델 (변경 없음)
- 접두사 있음 (`session.Session.Get`) → 패키지 Go interface 기반 검증

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `parser/types.go` | `Sequence`에 `Package` 필드 추가 |
| `parser/parser.go` | `parseCallExprInputs()` 에서 3-part 모델명 파싱 (`pkg.Model.Method`) |
| `parser/parser_test.go` | 패키지 접두사 파싱 테스트 |
| `validator/validator.go` | 내부 검증: 3-part 모델명 형식 허용. 외부 검증: 패키지 모델은 DDL 스킵 → 패키지 interface 검증 |
| `validator/symbol.go` | `LoadPackageInterfaces()` 추가: import 경로에서 Go interface 파싱 → `SymbolTable.Models`에 등록 |
| `validator/validator_test.go` | 패키지 모델 검증 테스트 |
| `generator/go_target.go` | 코드젠: 패키지 접두사 모델 호출 생성 (`sessionModel.Get(...)`) |
| `generator/generator_test.go` | 패키지 모델 코드젠 테스트 |

## 설계

### 1. 파서 — 3-part 모델명 파싱

현재 `parseCallExprInputs()`는 `Model.Method`를 통째로 `seq.Model`에 저장한다. 3-part 이름 지원:

```
session.Session.Get({...})
→ Package: "session", Model: "Session.Get" (기존 Model 필드 유지)
→ 또는 Package: "session", Model: "Session", Method: "Get" (Model 분리)
```

**방안**: `Sequence.Package` 필드 추가. `parseCallExprInputs()` 에서 dot이 2개 이상이면 첫 번째를 Package로 분리.

```go
// parser/types.go
type Sequence struct {
    // ...
    Package string // "session", "" (패키지 접두사)
    Model   string // "Session.Get" (기존과 동일 형식)
    // ...
}
```

파싱 규칙:
- `User.FindByID(...)` → Package="", Model="User.FindByID"
- `session.Session.Get(...)` → Package="session", Model="Session.Get"

### 2. 내부 검증 — 형식 검사

`validateRequiredFields()`의 Model 형식 검사:
- 현재: `SplitN(seq.Model, ".", 2)` → 2-part 확인
- 변경: Package 있으면 Model은 여전히 `Model.Method` 2-part이므로 기존 로직 유지

### 3. 외부 검증 — 패키지 모델 분기

`validateModel()`에서 `seq.Package != ""`일 때:
- DDL 모델 검증 스킵 (기존 `st.Models` 조회 대신)
- 패키지 interface 검증 경로로 분기

```go
if seq.Package != "" {
    // 패키지 모델: import 경로에서 interface 탐색
    pkgModelKey := seq.Package + "." + modelName // "session.Session"
    ms, ok := st.Models[pkgModelKey]
    if !ok {
        // WARNING: 검증 불가 — interface 미제공
    }
    // 메서드 존재 확인
} else {
    // 기존 DDL 모델 검증
}
```

### 4. 심볼 테이블 — 패키지 interface 로딩

`SymbolTable`에 패키지 interface 로딩 기능 추가. 서비스 파일의 import 경로를 활용:

```go
// symbol.go
func (st *SymbolTable) LoadPackageInterfaces(imports []string, projectRoot string) error
```

로딩 전략:
1. 서비스 파일의 import 선언에서 패키지명 → import 경로 매핑
2. import 경로가 상대 경로면 `projectRoot` 기준 탐색
3. Go interface 파싱 → `st.Models["session.Session"]` 형태로 등록
4. interface 없으면 WARNING (ERROR 아님)
5. 패키지 경로 자체가 없으면 ERROR

에러 메시지 (자기 교정 루프):
```
[ERROR] session.Session — 메서드 "Gett" 없음. 사용 가능: Set, Get, Delete
```

### 5. 코드젠 — 패키지 모델 호출

`go_target.go`에서 `seq.Package != ""`일 때 코드젠 분기:

```go
// 현재 (접두사 없음)
d.ModelCall = lcFirst(parts[0]) + "Model." + parts[1]
// → courseModel.FindByID(...)

// 패키지 접두사 있음
d.ModelCall = lcFirst(parts[0]) + "Model." + parts[1]
// → sessionModel.Get(...)
// (패키지명은 handler 레벨에서 DI 주입, 코드젠 형태는 동일)
```

모델 인터페이스 파생 (`renderInterfaces()`):
- 패키지 모델은 interface를 자체 제공하므로 `models_gen.go`에 생성하지 않음
- `renderInterfaces()`에서 `Package != ""` 모델 스킵

import 처리:
- 패키지 모델은 해당 패키지 import 추가 (`collectImports()`에서 처리)

### 6. 영향 없는 부분

- 접두사 없는 기존 모델 (`User.FindByID`) → 변경 없음
- @call (`auth.VerifyPassword`) → 기존 로직 유지 (SeqCall은 이미 별도 분기)
- @response, @empty, @exists, @state, @auth → Model 필드 미사용, 영향 없음

## 의존성

- fullend Phase 021 (패키지 접두사 모델 판단 규칙 + Func 순수성 강제)
- 패키지 경로 탐색: spec 파일 import 선언 활용 (기존 @call과 동일 패턴)

## 검증 방법

```bash
go test ./parser/... ./validator/... ./generator/... -count=1
```

테스트 케이스:
1. `TestParsePackagePrefixModel` — `session.Session.Get(...)` → Package="session", Model="Session.Get"
2. `TestParseNoPackagePrefix` — `User.FindByID(...)` → Package="", Model="User.FindByID" (기존 동작 유지)
3. `TestValidatePackageModelMethodExists` — 패키지 interface에 메서드 있음 → OK
4. `TestValidatePackageModelMethodMissing` — 메서드 없음 → ERROR (사용 가능 메서드 안내)
5. `TestValidatePackageModelNoInterface` — interface 없음 → WARNING
6. `TestValidatePackageModelSkipDDL` — 패키지 모델은 DDL 검증 스킵
7. `TestGeneratePackageModelCall` — `sessionModel.Get(...)` 코드젠
8. `TestGeneratePackageModelSkipInterface` — 패키지 모델은 `models_gen.go`에 생성 안함
