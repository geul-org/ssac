# SSaC — Service Sequences as Code

> Service 계층의 실행 흐름을 선언적 블록(sequence)으로 분해하여, 비즈니스 로직의 SSOT를 구성한다.

## 빈 공간

기존 SSOT는 함수 외부까지만 커버한다:

| 기존 SSOT | 커버 범위 | 함수 내부? |
|---|---|---|
| OpenAPI | API 경로, 파라미터, 응답 스키마 (→ Controller) | X |
| SQL DDL | 테이블 구조, 인덱스, 제약 (→ Model) | X |

**함수 내부 — "조회 → 생성 → 응답"이라는 비즈니스 흐름은 선언할 곳이 없다.** 구현 코드를 읽어야만 알 수 있다. 이 빈 공간을 sequence가 채운다. (입력 검증은 OpenAPI가 담당하므로 Controller 책임.)

## 타입 참조

SSaC는 자체 struct를 정의하지 않는다. sequence에서 사용하는 타입은 각 모델 SSOT에서 산출된 것을 참조한다.

| 모델 SSOT | 코드젠 | 산출 타입 | 예시 |
|---|---|---|---|
| SQL DDL (`specs/db/`) | sqlc | DB 모델 struct | `Project`, `Session` |
| OpenAPI schemas (`specs/api/`) | openapi-generator | 요청/응답 struct | `CreateSessionRequest` |
| Go interface (`specs/model/`) | — (선언 자체가 코드) | 비-DB 리소스 계약 | `FileSystem`, `Cache` |

sequence에서 `@result project Project`라고 쓰면, `Project`는 sqlc가 DDL에서 산출한 struct다. 심볼릭 검증 시 각 모델 SSOT의 타입 정의와 sequence의 변수 참조를 교차 체크한다.

변수와 타입은 Go 관례로 구분한다:
- `project` — 변수 (소문자)
- `Project` — 타입 (대문자)
- `ProjectID` — request 필드 (`request` 키워드와 함께 사용)
- `currentUser` — 예약어 (인증 컨텍스트에서 자동 주입, 선언 불필요)
- `config` — 예약어 (환경변수/ConfigMap에서 자동 주입, 선언 불필요)

## sequence란

함수 내부의 실행 블록을 타입화한 선언 단위.

**what(뭘 하는가)만 선언하고, how(어떻게 하는가)는 코드젠이 채운다.**

### @model의 범위

`@model`은 DB 테이블에 한정되지 않는다. CRUD로 다룰 수 있는 모든 리소스(파일시스템, 외부 API, 캐시 등)가 Model이다. `Project.FindByID`가 DB 조회일 수도, `FileSystem.Exists`가 파일 존재 확인일 수도, `GitHub.CreateRepo`가 외부 API 호출일 수도 있다. sequence의 get/post/put/delete는 리소스 종류를 가리지 않는다.

### @message — 모든 타입의 공통 선택 필드

모든 sequence 타입은 선택적으로 `@message`를 가질 수 있다. 기본 메시지는 타입과 모델명 조합에서 자동 생성되며, 커스텀이 필요할 때만 선언한다.

자동 생성 규칙:
- `get` + `Project.FindByID` → "Project 조회 실패"
- `post` + `Session.Create` → "Session 생성 실패"
- `guard nil` + `project` → "project가 존재하지 않습니다"
- `delete` + `Project.Delete` → "Project 삭제 실패"

커스텀 예시:
```go
// @sequence get
// @model Project.FindByID
// @param ProjectID request
// @result project Project
// @message "해당 프로젝트를 찾을 수 없습니다"
```

`@message`가 없으면 코드젠이 기본값을 사용한다. 기존 guard의 `@error`는 `@message`로 통합된다.

### sequence 타입

| 타입 | 역할 | 예시 |
|---|---|---|
| authorize | 권한 검증 | @action delete @resource project @id ProjectID → OPA 질의 |
| get | 리소스 조회 | Model.FindByID(id) → result Type |
| guard nil | 결과가 null이면 종료 | result가 없으면 메시지 후 종료 |
| guard exists | 결과가 있으면 종료 | result가 이미 존재하면 메시지 후 종료 |
| post | 리소스 생성 | Model.Create(fields...) → result Type |
| put | 리소스 수정 | Model.Update(id, fields...) |
| delete | 리소스 삭제 | Model.Delete(id) |
| password | 비밀번호 비교 | hash 비교 후 실패 시 종료 |
| call | 외부 호출 | @component → 등록된 component (코드젠 자동), @func → 순수 함수 (직접 구현) |
| response | 응답 반환 | json, view, redirect |

### 표현 예시 — 단순 (생성)

```go
// specs/backend/service/create_session.go

// @sequence get
// @model Project.FindByID
// @param ProjectID request
// @result project Project

// @sequence guard nil project
// @message "프로젝트가 존재하지 않습니다"

// @sequence post
// @model Session.Create
// @param ProjectID request
// @param Command   request
// @result session Session

// @sequence response json
// @var session
func CreateSession(w http.ResponseWriter, r *http.Request) {}
```

### 표현 예시 — 복합 (삭제 + 권한 + 비즈니스 검증 + 별도 함수)

```go
// specs/backend/service/delete_project.go

// @sequence authorize
// @action delete
// @resource project
// @id ProjectID

// @sequence get
// @model Project.FindByID
// @param ProjectID request
// @result project Project

// @sequence guard nil project
// @message "프로젝트가 존재하지 않습니다"

// @sequence get
// @model Session.CountByProjectID
// @param ProjectID request
// @result sessionCount int

// @sequence guard exists sessionCount
// @message "하위 세션이 존재하여 삭제할 수 없습니다"

// @sequence call
// @component notification
// @param project.OwnerEmail
// @param "프로젝트가 삭제됩니다"

// @sequence call
// @func cleanupProjectResources
// @param project
// @result cleaned bool

// @sequence delete
// @model Project.Delete
// @param ProjectID request

// @sequence response json
func DeleteProject(w http.ResponseWriter, r *http.Request) {}
```

### 코드젠 결과 (파생되는 구현)

```go
// artifacts/backend/internal/service/create_session.go

func CreateSession(w http.ResponseWriter, r *http.Request) {
    // 입력 검증은 Controller(OpenAPI 코드젠)에서 완료됨

    // get
    project, err := projectModel.FindByID(projectID)
    if err != nil {
        http.Error(w, "Project 조회 실패", http.StatusInternalServerError)
        return
    }

    // guard nil
    if project == nil {
        http.Error(w, "프로젝트가 존재하지 않습니다", http.StatusNotFound)
        return
    }

    // post
    session, err := sessionModel.Create(projectID, command)
    if err != nil {
        http.Error(w, "Session 생성 실패", http.StatusInternalServerError)
        return
    }

    // response json
    json.NewEncoder(w).Encode(map[string]interface{}{
        "session": session,
    })
}
```

## 정합성 검증

SSaC는 서비스 로직 선언일 뿐 아니라, OpenAPI와 DDL 사이의 **정합성 검증 허브**다. sequence가 양쪽 SSOT를 동시에 참조하므로, SSaC 검증이 통과하면 세 SSOT 간 정합성이 보장된다.

```
OpenAPI (specs/api/)
    ↘
      SSaC (specs/backend/service/)  ← 교차 검증
    ↗
DDL (specs/db/)
```

### OpenAPI ↔ SSaC

- `@param ProjectID request` → 해당 엔드포인트의 request body 또는 path parameter에 `ProjectID`가 존재하는가?
- `@sequence response json` + `@var session` → 해당 엔드포인트의 response schema에 `session` 타입이 일치하는가?

### DDL ↔ SSaC

- `@model Project.FindByID` → DDL에 `projects` 테이블이 존재하고, sqlc query에 `FindByID`가 정의되어 있는가?
- `@param ProjectID request` → 해당 쿼리의 파라미터 타입과 일치하는가?
- `@result project Project` → sqlc가 산출할 struct와 매칭되는가?

### OpenAPI ↔ DDL (간접 보장)

SSaC가 OpenAPI의 request 필드를 받아서 DDL의 쿼리에 넘기므로, SSaC 검증이 통과하면 OpenAPI와 DDL 간 타입 정합성도 자동으로 보장된다. 별도의 OpenAPI ↔ DDL 직접 검증이 필요 없다.

### 검증 파이프라인

```
1. sqlc generate          → Go struct 산출
2. openapi-generator      → Go struct 산출
3. SSaC validate          → 1, 2의 심볼 테이블 수집 → sequence 참조 교차 체크
4. SSaC codegen           → 검증 통과한 sequence만 코드 생성
5. go build               → 최종 확인
```

3단계에서 OpenAPI spec(yaml/json)과 sqlc 설정에서 쿼리 목록을 읽어 심볼 테이블을 구성한다. sequence의 모든 `@model`, `@param`, `@result` 참조를 이 테이블에서 lookup한다. 불일치가 있으면 코드 생성 전에 에러를 보고한다.

## 설계 원칙

### early return 패턴

sequence가 1차원 나열로 성립하는 전제는 **early return(guard clause) 패턴**이다. 모든 검증 sequence(guard nil, guard exists, authorize)는 `if 조건 { return }` 형태로 실패 시 즉시 탈출한다. 분기 없이 위에서 아래로 읽으면 끝이다.

```
authorize → if !allowed { return }     ← guard
get       → 조회
guard nil → if nil { return }          ← guard
post      → 실행
response  → 반환
```

if/else로 중첩하면 depth가 깊어져 코드가 꼬이고, sequence의 선형 나열이 불가능해진다. Go 공식 문서도 이 패턴을 권장한다 — "indent error flow, not the happy path."

### sequence로 표현 못하는 로직은 call로 위임한다

```
Service 함수
  ├─ sequence 기본 타입  → 프레임워크 레벨 (고정, 모든 프로젝트 공통, 심볼릭 코드젠)
  ├─ call @component    → 프로젝트 레벨 (등록제, 프로젝트별 확장, 심볼릭 코드젠)
  └─ call @func         → 비즈니스 레벨 (고유 로직, LLM/사람이 직접 구현)
```

분기 후 합류, 루프 내 조건부 처리, 트랜잭션 보상 등 복잡한 로직은 sequence를 확장하지 않고 call로 위임한다. 반복 패턴이 3번 이상 나타나면 component로 승격한다. component는 `specs/backend/component/`에 정의하며 심볼릭 코드젠이 가능하다. 순수 함수(@func)는 테스트 가능한 단위로 격리된다.

## SSOT 계층 구조

```
OpenAPI          → Controller  "어떤 API가 있는가"        (라우팅, 파라미터 바인딩) — openapi-generator
Model SSOT       → Model       "어떤 리소스가 있는가"      (CRUD, 계약)
  ├─ SQL DDL     →   DB 모델    (테이블, 제약, 인덱스)     — sqlc
  ├─ OpenAPI     →   외부 API   (요청/응답 스키마)         — openapi-generator
  └─ Go interface →  비-DB 모델  (파일시스템, 캐시 등)      — 선언 자체가 코드
SSaC             → Service     "어떤 흐름으로 처리하는가"   (get → post → response) — SSaC codegen
                   └─ 구현 코드  "어떻게 하는가"            (if err != nil, sql.QueryRow, ...)
```

위에서 아래로 갈수록 구체적이고, SSOT는 위 계층만 선언한다. 최하층(구현 코드)은 코드젠이 산출한다.

## 토큰 비용 비교

| | SSOT (sequence) | 구현 코드 |
|---|---|---|
| 줄 수 | 10~15줄 | 30~100줄 |
| 표현 내용 | 의도 (what) | 의도 + 구현 (what + how) |
| 에러 핸들링 | 없음 | 있음 |
| 라이브러리 의존 | 없음 | http, sql, json 등 |
| 변경 시점 | 요구사항 변경 시 | 리팩토링, 버그 수정 포함 |

## 심볼릭 코드젠 가능성

sequence는 **타입이 고정**되어 있다 (authorize, get, guard nil, guard exists, post, put, delete, password, call, response). 각 타입의 코드젠 템플릿을 Go로 작성하면 LLM 없이 심볼릭으로 구현 코드를 산출할 수 있다.

```
sequence 파싱 (Go AST / 코멘트 파싱)
  → 타입별 템플릿 매칭
    → Go 코드 생성
```

이것이 SSaC의 심볼릭 코드젠을 `과도기: LLM`에서 자동화로 전환할 수 있는 근거다.
