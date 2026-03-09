✅ 완료

# Phase 8: 도메인 폴더 구조 지원

수정지시서 007 기반. 서비스 파일을 도메인별 폴더로 분류할 수 있도록 parser와 generator를 확장한다.

## 목표

- `specs/service/auth/login.go` 같은 도메인 폴더 구조 지원
- 기존 flat 구조(`specs/service/login.go`) 하위 호환 유지
- 생성 시 도메인별 서브디렉토리 + 해당 package 이름 사용

## 작업 순서

### Step 1: ServiceFunc에 Domain 필드 추가

`parser/types.go`:
```go
type ServiceFunc struct {
    Name      string
    FileName  string
    Domain    string     // 도메인 폴더명 (e.g. "auth"). 빈 문자열이면 루트.
    Sequences []Sequence
}
```

### Step 2: ParseDir 재귀 탐색 + Domain 파생

`parser/parser.go` — `os.ReadDir` → `filepath.WalkDir`로 변경:

```go
func ParseDir(dir string) ([]ServiceFunc, error) {
    var funcs []ServiceFunc
    err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
        if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".go") {
            return err
        }
        sf, err := ParseFile(path)
        if err != nil {
            return fmt.Errorf("%s 파싱 실패: %w", path, err)
        }
        if sf != nil {
            rel, _ := filepath.Rel(dir, path)
            if parts := strings.Split(filepath.Dir(rel), string(filepath.Separator)); parts[0] != "." {
                sf.Domain = parts[0]
            }
            funcs = append(funcs, *sf)
        }
        return nil
    })
    return funcs, err
}
```

하위 호환:
- flat: `service/login.go` → `rel="login.go"`, `Dir(rel)="."` → `Domain=""`
- 폴더: `service/auth/login.go` → `rel="auth/login.go"`, `Dir(rel)="auth"` → `Domain="auth"`

### Step 3: GenerateWith 도메인별 출력 경로

`generator/generator.go` — `GenerateWith` 수정:

```go
for _, sf := range funcs {
    code, err := t.GenerateFunc(sf, st)
    // ...
    outPath := outDir
    if sf.Domain != "" {
        outPath = filepath.Join(outDir, sf.Domain)
        os.MkdirAll(outPath, 0755)
    }
    path := filepath.Join(outPath, outName)
    // ...
}
```

### Step 4: GenerateFunc 도메인별 package 이름

`generator/go_target.go` — 현재 `package service` 하드코딩 → Domain 기반:

```go
pkgName := "service"
if sf.Domain != "" {
    pkgName = sf.Domain
}
buf.WriteString(fmt.Sprintf("package %s\n\n", pkgName))
```

### Step 5: 테스트

기존 테스트 전부 통과 (flat 구조 하위 호환) + 신규 테스트:

1. `TestParseDirRecursive` (parser_test.go): testdata에 도메인 폴더 fixture 추가, Domain 필드 확인
2. `TestGenerateDomain` (generator_test.go): Domain 있는 ServiceFunc → package 이름 + 출력 경로 확인

testdata fixture:
```
testdata/domain-service/
├── login.go           ← flat, Domain=""
└── course/
    └── create_course.go  ← Domain="course"
```

## 변경 파일 목록

| 파일 | 변경 유형 |
|---|---|
| `parser/types.go` | 수정: `ServiceFunc`에 `Domain string` 추가 |
| `parser/parser.go` | 수정: `ParseDir` → `filepath.WalkDir` + Domain 파생, import에 `io/fs` 추가 |
| `parser/parser_test.go` | 추가: `TestParseDirRecursive` |
| `generator/generator.go` | 수정: `GenerateWith` 도메인별 출력 경로 |
| `generator/go_target.go` | 수정: `GenerateFunc` 도메인별 package 이름 |
| `generator/generator_test.go` | 추가: `TestGenerateDomain` |
| `testdata/domain-service/` | 신규: 도메인 폴더 테스트 fixture |
| `testdata/domain-service/course/` | 신규: 도메인 서브디렉토리 fixture |

## 하지 않는 것

- 모델 인터페이스 도메인 분리 — 모델은 도메인 간 공유되므로 단일 `model/` 유지 (지시서 §4)

## 검증

```bash
go test ./parser/... ./generator/... -count=1

# flat 구조: Domain="" → outDir/login.go, package service
# 폴더 구조: Domain="course" → outDir/course/create_course.go, package course
```

## 리스크

- 낮음: 순수 확장, 기존 기능 변경 없음 (Domain="" 분기)
- `filepath.WalkDir`은 알파벳 순 → 기존 `os.ReadDir`와 동일 정렬
