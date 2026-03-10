✅ 완료

# Phase 015: 모델 필드 참조 시 포인터 역참조 코드 생성

## 목표

모델 result 변수의 필드를 참조할 때, 생성 코드에서 항상 `*source.Field` 역참조를 생성한다. 모든 모델 필드를 포인터(`*type`)로 통일 취급하여 일관성 확보.

```go
// 현재 생성 코드
billing.ReleaseFunds(billing.ReleaseFundsRequest{FreelancerID: gig.FreelancerID})
reservationModel.Create(endAt, roomID, startAt, currentUser.ID)

// 기대 생성 코드
billing.ReleaseFunds(billing.ReleaseFundsRequest{FreelancerID: *gig.FreelancerID})
reservationModel.Create(*reservation.EndAt, roomID, *reservation.StartAt, currentUser.ID)
```

## 설계

NOT NULL 여부를 구분하지 않고 모델 필드는 전부 포인터로 취급한다. NOT NULL 검증은 DB가 담당.

### 역참조 대상

`source.Field` 형식에서 source가 **모델 result 변수**일 때만 `*` 접두사:

| source | 역참조 | 이유 |
|---|---|---|
| `reservation.Status` | `*reservation.Status` | 모델 result → 필드가 포인터 |
| `request.CourseID` | `courseID` (변환) | HTTP 파라미터 → 이미 값 타입 |
| `currentUser.ID` | `currentUser.ID` | 인증 컨텍스트 → 값 타입 |
| `config.Secret` | `config.Secret` | 설정 → 값 타입 |
| `"cancelled"` | `"cancelled"` | 리터럴 → 값 |

예약 소스(`request`, `currentUser`, `config`, `query`)가 아닌 `source.Field`가 역참조 대상.

### 구현

`inputValueToCode()`에서 모델 변수 필드 참조를 감지하고 `*` 접두사:

```go
func inputValueToCode(val string) string {
    if val == "query" { return "opts" }
    if strings.HasPrefix(val, "request.") { return lcFirst(val[len("request."):]) }
    if strings.HasPrefix(val, "currentUser.") || strings.HasPrefix(val, "config.") {
        return val
    }
    // 리터럴
    if strings.HasPrefix(val, `"`) { return val }
    // 모델 result 변수 필드 → 역참조
    if strings.ContainsRune(val, '.') {
        return "*" + val
    }
    // bare variable (hashedPassword 등)
    return val
}
```

DDLTable, SymbolTable 변경 없음. 함수 시그니처 변경 없음.

## 변경 파일

| 파일 | 내용 |
|---|---|
| `generator/go_target.go` | `inputValueToCode()` — 모델 필드 참조 시 `*` 접두사 |
| `generator/generator_test.go` | 역참조 테스트 추가, 기존 테스트 expected output 업데이트 |

## 검증

```bash
go test ./parser/... ./validator/... ./generator/... -count=1
ssac gen specs/dummy-study/ /tmp/ssac-phase15-check/
```

## 의존성

- 수정지시서v2/010
