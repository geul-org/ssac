# ssac

## 프로젝트 루트
~/.clari/repos/ssac

## 프로젝트 개요
Service Sequences as Code — Go 주석 기반 선언적 서비스 로직을 파싱하여 Go 구현 코드를 생성하는 CLI 도구.

## CLI

```
ssac parse [dir]              # 주석 파싱 결과 출력 (기본: specs/backend/service/)
ssac validate [dir]           # 내부 검증 또는 외부 SSOT 교차 검증 (자동 감지)
ssac gen <service-dir> <out>  # validate → codegen → gofmt (심볼 테이블 있으면 타입 변환 + 모델 인터페이스 생성)
```

## 계획 작성 원칙

구현 전 `specs/plans/`에 계획 md를 작성한다.

- 파일명: `PhaseNNN-TITLE.md` (예: `Phase001-CLISkeleton.md`)
- 구현 코드를 쓰기 전에 계획을 먼저 작성하고 승인을 받는다
- 계획에는 다음을 포함한다:
  - 목표: 무엇을 만드는가
  - 변경 파일 목록: 어떤 파일을 생성/수정하는가
  - 의존성: 외부 패키지, 형제 프로젝트 API
  - 검증 방법: 어떻게 확인하는가
- 계획이 승인되면 구현하고, 완료 후 계획 상단에 `✅ 완료` 표시

## 기술 스택

- Go 1.24+, module: `github.com/geul-org/ssac`
- 파싱: `go/ast`, `go/parser`
- 코드젠: `text/template`, `go/format`
- 외부 의존성: `gopkg.in/yaml.v3` (OpenAPI 파싱)

## DSL 문법

```go
// @sequence <type>        — 블록 시작. 10종: authorize|get|guard nil|guard exists|guard state|post|put|delete|call|response
// @model <Model.Method>   — 리소스 모델.메서드 (get/post/put/delete)
// @param <Name> <source>  — source: request, currentUser, 변수명, "리터럴"
// @result <var> <Type>    — 결과 바인딩 (get/post 필수, call 선택)
// @message "msg"          — 커스텀 에러 메시지 (선택, 기본값 자동생성)
// @var <name>             — response에서 반환할 변수
// @action @resource @id   — authorize 전용 (3개 모두 필수)
// @component | @func      — call 전용 (택일 필수). @func는 package.funcName 형식 필수
```

타입별 필수 태그:

| 타입 | 필수 |
|---|---|
| authorize | @action, @resource, @id |
| get, post | @model, @result |
| put, delete | @model |
| guard nil/exists | target (sequence 라인에 변수명) |
| guard state | target (stateDiagramID), @param 1개 (entity.Field) |
| call | @component 또는 @func package.funcName (택일) |
| response | (없음, @var는 선택) |

## 디렉토리

```
cmd/ssac/main.go                 # CLI 진입점
parser/                          # 주석 → []ServiceFunc
validator/                       # 내부 + 외부 SSOT 검증
generator/                       # Target 인터페이스 기반 코드젠 (다중 언어 확장 가능)
  target.go                      #   Target 인터페이스 + DefaultTarget()
  go_target.go                   #   GoTarget: Go 코드 생성 구현
  go_templates.go                #   Go 템플릿
  generator.go                   #   하위 호환 래퍼 (Generate, GenerateWith) + 유틸
specs/                           # 선언 (입력, SSOT)
  dummy-study/                   #   스터디룸 예약 더미 프로젝트
    service/  db/queries/  api/  model/
  plans/                         #   구현 계획서
artifacts/                       # 문서
  manual-for-human.md            #   상세 매뉴얼 (인간용)
  manual-for-ai.md               #   컴팩트 레퍼런스 (AI용)
testdata/                        # 테스트 fixture
files/                           # 기초 자료
  SSaC.md                        #   기획서
```

## 외부 검증 프로젝트 구조

`ssac validate <project-root>` 시 자동 감지:
- `<root>/service/**/*.go` — sequence spec (재귀 탐색, 도메인 폴더 지원)
- `<root>/db/*.sql` — DDL (CREATE TABLE → 컬럼 타입)
- `<root>/db/queries/*.sql` — sqlc 쿼리 (파일명→모델, `-- name: Method :cardinality`)
- `<root>/api/openapi.yaml` — OpenAPI 3.0 (operationId=함수명, x-pagination/sort/filter/include)
- `<root>/model/*.go` — Go interface→component. @func는 외부 패키지이므로 교차검증 스킵

## 코드젠 기능

심볼 테이블(외부 SSOT)이 있을 때 추가되는 기능:

- **타입 변환 코드젠**: DDL 컬럼 타입 기반으로 request 파라미터 변환 코드 생성 (int64→`strconv.ParseInt`, time.Time→`time.Parse`, 실패 시 400 early return)
- **Guard 값 타입**: result 타입에 따른 zero value 비교 (int→`== 0`/`> 0`, string→`== ""`/`!= ""`, pointer→`== nil`/`!= nil`)
- **currentUser/config source**: `@param Name currentUser` → `currentUser.Name`
- **Stale 데이터 경고**: put/delete 후 갱신 없이 response에 사용하면 WARNING
- **QueryOpts 자동 전달**: x-확장 있으면 `opts := QueryOpts{}` 생성 + 모델 호출에 `opts` 인자 자동 추가
- **List 3-tuple 반환**: many + QueryOpts → `result, total, err :=` (count 포함)
- **모델 인터페이스 파생**: 3 SSOT 교차(sqlc 카디널리티 + SSaC 파라미터 + OpenAPI x-확장) → `<outDir>/model/models_gen.go`
  - 모든 @param 소스 포함: request, currentUser, dot notation(`user.ID` → `userID`), 리터럴(`"pending"` → DDL 역매핑)
- **도메인 폴더 구조**: `service/auth/login.go` → `Domain="auth"` → `outDir/auth/login.go`, `package auth`
  - flat 구조(`service/login.go`) 하위 호환 유지 (Domain="")
- **@func 패키지 코드젠**: `@func auth.verifyPassword` → `auth.VerifyPassword(auth.VerifyPasswordInput{...})` 호출 코드 생성
  - `@result` 없으면 guard형 (401), 있으면 value형 (500)

## 더미 프로젝트

ssac 자체 더미: `specs/dummy-study/` (내부 테스트용, 외부 검증 프로젝트 구조)

fullend 더미 (SSaC 소비자, 통합 검증용):

| 프로젝트 | SSOT (specs) | 생성 산출물 (artifacts) |
|---|---|---|
| dummy-study | `~/.clari/repos/fullend/specs/dummy-study/` | `~/.clari/repos/fullend/artifacts/dummy-study/` |
| dummy-lesson | `~/.clari/repos/fullend/specs/dummy-lesson/` | `~/.clari/repos/fullend/artifacts/dummy-lesson/` |

각 더미 프로젝트 구조:
- `specs/<project>/frontend/*.html` — STML 페이지
- `specs/<project>/frontend/components/*.tsx` — 커스텀 컴포넌트
- `specs/<project>/api/openapi.yaml` — OpenAPI 스펙
- `artifacts/<project>/frontend/src/*.tsx` — 생성된 React 페이지
- `artifacts/<project>/backend/` — 생성된 Go 백엔드

## Coding Conventions

- gofmt 준수, 에러 즉시 처리 (early return)
- 파일명: snake_case, 변수/함수: camelCase, 타입: PascalCase
- 테스트: `go test ./parser/... ./validator/... ./generator/... -count=1`
- 테스트용 fixture는 `testdata/`에 배치. `/tmp` 등 외부 경로 사용 금지.

## Git 규칙

- Co-Authored-By 금지. 커밋 메시지에 AI 이름을 절대 포함하지 않는다.
- remote: `https://github.com/geul-org/ssac.git`
- 라이선스: MIT
