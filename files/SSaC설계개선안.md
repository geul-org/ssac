# SSaC 설계 개선안

더미 스터디 프로젝트에 `ssac gen`을 적용한 결과에서 도출한 개선 사항.

## 1. 버그: guard exists + 값 타입 → 컴파일 에러

### 현상

```go
// spec
// @result reservationCount int
// @sequence guard exists reservationCount

// 생성 코드
if reservationCount != nil {  // ← int는 nil 비교 불가
```

### 원인

`guard exists` 템플릿이 `{{.Target}} != nil`을 하드코딩한다.
result 타입이 포인터(`*Reservation`)면 정상이지만, 값 타입(`int`, `bool`)이면 컴파일 에러.

### 해결

Result.Type을 참조하여 guard 조건을 분기한다.

| Result.Type | guard nil 조건 | guard exists 조건 |
|---|---|---|
| 포인터/struct (기본) | `== nil` | `!= nil` |
| `int`, `int64` | `== 0` | `> 0` |
| `bool` | `== false` | `== true` |
| `string` | `== ""` | `!= ""` |

**구현 방안**: generator의 `buildTemplateData`에서 직전 sequence의 result 타입을 조회하고, 타입에 맞는 zero value 비교를 `templateData`에 전달한다. 또는 guard용 템플릿을 타입별로 분리한다.

```go
// templateData에 추가
ZeroCheck string  // "== nil", "== 0", "== false" 등
ExistsCheck string // "!= nil", "> 0", "== true" 등
```

파서가 이미 Result.Type을 파싱하므로 (`int`, `[]Reservation` 등) generator만 수정하면 된다.


## 2. 버그: currentUser 파라미터 변수 미선언

### 현상

```go
// spec
// @param UserID currentUser

// 생성 코드
reservationModel.ListByUserID(userID)  // ← userID 어디서도 선언 안 됨
```

`currentUser` source 파라미터는 `r.FormValue`로 추출되지 않는다 (정상).
그러나 변수 바인딩 코드도 생성되지 않아 미선언 변수를 참조한다.

### 원인

`collectRequestParams`는 source가 `"request"`인 것만 수집한다. source가 `"currentUser"`인 파라미터는 어디서도 변수를 생성하지 않는다.

### 해결

SSaC 기획서 정의: `currentUser`는 **예약어로, 인증 컨텍스트에서 자동 주입**된다.

두 가지 방향:

**(A) 코드젠에서 예약어 변수 추출 코드를 생성**

```go
// source가 "currentUser"인 파라미터가 있으면 함수 상단에 삽입
userID := currentUser.ID  // 또는 context에서 추출
```

문제: `currentUser`의 구체적 타입과 필드 접근 방식이 프로젝트마다 다를 수 있다.

**(B) 예약어는 이미 스코프에 있다고 가정하고, resolveParamRef에서 source를 활용**

현재 `resolveParamRef`는 Name만 보고 변환한다. source가 `"currentUser"`이면 `currentUser.{Name}` 형태로 변환:

```go
// @param UserID currentUser → currentUser.UserID
// @param ID currentUser     → currentUser.ID
```

이 방식이면 `currentUser`가 함수 인자나 context로 주입되어 있다는 전제 하에 작동한다. 함수 시그니처 자체의 개선(아래 4번)과 연결된다.

**권장**: (B) 방식. `resolveParamRef`에서 source가 예약어(`currentUser`, `config`)이면 `source.Name` 형태로 변환.


## 3. 설계 개선: Model interface 파생 생성

### 현상

생성 코드에서 `userModel.FindByEmail(email)`, `reservationModel.ListByUserID(userID)` 등을 호출하지만, 이 모델 객체의 **interface 정의가 존재하지 않는다**. sqlc는 `Queries.FindByEmail(ctx, email string)` 형태의 concrete 구현을 생성하고, SSaC는 `userModel.FindByEmail(email)` 형태로 호출한다. 이 둘 사이를 연결하는 계약이 부재하다.

### 현재 SSOT 구조의 빈 자리

```
OpenAPI  → Controller  "어떤 API가 있는가"       ← spec 있음
???      → Model       "어떤 리소스 계약이 있는가"  ← spec 없음 ★
SSaC     → Service     "어떤 흐름으로 처리하는가"   ← spec 있음
SQL DDL  → DB          "어떤 테이블이 있는가"      ← spec 있음
```

### 해결: 별도 spec 없이 기존 SSOT 교차에서 파생

Model interface를 정의하는 새로운 SSOT(MaC)를 만들 필요 없다. 필요한 정보가 이미 세 곳에 분산되어 있으며, 이를 교차하면 interface가 도출된다.

```
sqlc 쿼리    →  메서드명, 반환 카디널리티, 기본 파라미터 타입
SSaC spec   →  비즈니스 파라미터명과 출처
OpenAPI x-  →  인프라 파라미터 (pagination, sort, filter)
```

#### 교차 예시 1: 단건 조회

```
sqlc:       -- name: FindByEmail :one
            SELECT * FROM users WHERE email = $1;

SSaC:       @model User.FindByEmail
            @param Email request
            @result user User

파생:       type UserModel interface {
                FindByEmail(email string) (*User, error)
            }
```

#### 교차 예시 2: 목록 조회 + 인프라 파라미터

```
sqlc:       -- name: ListByUserID :many
            SELECT * FROM reservations WHERE user_id = $1;

SSaC:       @model Reservation.ListByUserID
            @param UserID currentUser
            @result reservations []Reservation

OpenAPI:    x-pagination: { defaultLimit: 20, maxLimit: 100 }
            x-sort: { allowed: [start_at, created_at] }

파생:       type ReservationModel interface {
                ListByUserID(userID string, opts QueryOpts) ([]Reservation, int, error)
            }
```

#### 각 SSOT의 기여

| 관심사 | 선언하는 곳 | interface에 기여 |
|---|---|---|
| 메서드 존재 여부, 반환 타입 | sqlc 쿼리 | 메서드명, `:one`→포인터 / `:many`→슬라이스 / `:exec`→error만 |
| 비즈니스 파라미터 | SSaC spec | `userID string` 등 호출 인자 |
| 페이지네이션/정렬/필터 | OpenAPI x- | `opts QueryOpts` (해당 operation에 x-가 있을 때만) |
| 파라미터 Go 타입 | sqlc DDL 컬럼 타입 | `string`, `int64`, `time.Time` 등 |

### 구현 방안

`ssac gen` 실행 시 심볼 테이블 교차 단계에서 Model interface를 파생 산출물로 함께 생성한다.

1. **심볼 테이블 확장**: sqlc 쿼리의 카디널리티(`:one`/`:many`/`:exec`)와 DDL 컬럼 타입을 `ModelSymbol`에 저장
2. **교차 합성**: SSaC spec의 `@model` 사용 패턴 + OpenAPI `x-` 정보를 결합하여 메서드 시그니처 결정
3. **interface 코드젠**: 합성된 시그니처로 Go interface 파일 생성
4. **wrapper 코드젠**: sqlc `Queries` struct를 래핑하여 생성된 interface를 구현하는 adapter 생성


## 4. 설계 개선: 타입 변환

### 현상

모든 request 파라미터가 `r.FormValue`로 추출되어 `string` 타입이다.

```go
capacity := r.FormValue("Capacity")  // string
roomModel.Update(roomID, name, capacity, location)  // capacity는 int여야 할 수 있음
```

### 원인

파서가 파라미터의 **Go 타입 정보**를 갖고 있지 않다. `@param Name request`에는 타입 선언이 없다.

### 해결 방향

**(A) 외부 SSOT에서 타입 추론**

OpenAPI schema나 sqlc 쿼리 파라미터에서 타입을 가져온다. `ssac gen`이 심볼 테이블을 활용하여 request 파라미터의 타입을 결정하고, `strconv.Atoi` 등 변환 코드를 생성.

**(B) @param에 타입 힌트 추가**

```go
// @param Capacity int request
```

DSL 문법을 확장하면 파서만으로 타입 정보를 가질 수 있지만, OpenAPI와 중복 선언이 된다. SSaC의 "SSOT 교차 검증" 철학과 충돌.

**권장**: (A) 방식. 심볼 테이블에 파라미터 타입 정보를 추가하고, gen 시 활용.


## 5. 설계 개선: put 후 stale 데이터 반환

### 현상

`UpdateRoom`에서 `Room.FindByID`로 가져온 `room`을 수정(`Room.Update`) 후 그대로 response에 반환한다. 클라이언트는 수정 전 데이터를 받는다.

### 원인

이건 spec 자체의 문제이기도 하지만, 코드젠 차원에서 경고할 수 있다.

### 해결

validator에 규칙 추가: `put`/`delete` 이후 response에서 같은 모델의 변수를 반환하면 경고.

```
WARNING: update_room.go:UpdateRoom — "room"이 put(Room.Update) 이후 갱신 없이 response에 사용됩니다
```

spec 작성자가 re-fetch를 추가하거나 의도적임을 확인하도록 유도.


## 6. 개선 우선순위

| 순위 | 항목 | 영향 | 난이도 |
|---|---|---|---|
| 1 | guard exists 값 타입 분기 | 컴파일 에러 | 낮음 |
| 2 | currentUser 파라미터 해결 | 컴파일 에러 | 낮음 |
| 3 | stale 데이터 경고 | 검증 품질 | 낮음 |
| 4 | Model interface 파생 생성 | 실용성 · 구조적 | 중간 |
| 5 | 타입 변환 코드젠 | 실용성 | 중간 (4번에 의존) |

1, 2번은 생성 코드가 컴파일되지 않는 버그이므로 즉시 수정 대상이다.
4번이 해결되면 의존성 주입, 타입 변환(5번), OpenAPI x- 인프라 파라미터 배달이 모두 같은 기반 위에서 처리된다.
