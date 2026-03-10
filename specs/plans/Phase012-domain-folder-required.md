✅ 완료

# Phase 012: service/ 직접 .go 파일 금지 — 도메인 서브 폴더 강제

## 목표

`service/` 직접에 SSaC 파일을 두면 에러를 반환한다. `service/<domain>/*.go` 서브 폴더 필수.

## 변경 파일

| 파일 | 내용 |
|---|---|
| `parser/parser.go` | `ParseDir()`: `filepath.Dir(rel) == "."` → 에러 반환 |
| `parser/parser_test.go` | `TestParseFlatServiceError` 추가 |
| `specs/dummy-study/service/` | 도메인 서브 폴더로 마이그레이션 (auth, reservation, room) |

## 검증

```bash
go test ./parser/... ./validator/... ./generator/... -count=1
ssac gen specs/dummy-study/ /tmp/check/
```

## 의존성

- 수정지시서v2/006
