package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/geul-org/ssac/parser"
	"github.com/geul-org/ssac/validator"
)

// Generate는 []ServiceFunc를 받아 outDir에 Go 파일을 생성한다.
func Generate(funcs []parser.ServiceFunc, outDir string, st *validator.SymbolTable) error {
	return GenerateWith(DefaultTarget(), funcs, outDir, st)
}

// GenerateFunc는 단일 ServiceFunc의 Go 코드를 생성한다.
func GenerateFunc(sf parser.ServiceFunc, st *validator.SymbolTable) ([]byte, error) {
	return DefaultTarget().GenerateFunc(sf, st)
}

// GenerateModelInterfaces는 심볼 테이블과 SSaC spec을 교차하여 Model interface를 생성한다.
func GenerateModelInterfaces(funcs []parser.ServiceFunc, st *validator.SymbolTable, outDir string) error {
	return DefaultTarget().GenerateModelInterfaces(funcs, st, outDir)
}

// GenerateWith는 지정된 Target으로 코드를 생성한다.
func GenerateWith(t Target, funcs []parser.ServiceFunc, outDir string, st *validator.SymbolTable) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("출력 디렉토리 생성 실패: %w", err)
	}

	for _, sf := range funcs {
		code, err := t.GenerateFunc(sf, st)
		if err != nil {
			return fmt.Errorf("%s 코드 생성 실패: %w", sf.Name, err)
		}

		ext := t.FileExtension()
		outName := strings.TrimSuffix(sf.FileName, ".go") + ext
		outPath := outDir
		if sf.Domain != "" {
			outPath = filepath.Join(outDir, sf.Domain)
			os.MkdirAll(outPath, 0755)
		}
		path := filepath.Join(outPath, outName)
		if err := os.WriteFile(path, code, 0644); err != nil {
			return fmt.Errorf("%s 파일 쓰기 실패: %w", path, err)
		}
	}
	return nil
}

// lcFirst는 첫 글자를 소문자로 변환한다.
func lcFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// ucFirst는 첫 글자를 대문자로 변환한다.
func ucFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// toSnakeCase는 PascalCase/camelCase를 snake_case로 변환한다.
func toSnakeCase(s string) string {
	var result []byte
	for i, c := range s {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				prev := s[i-1]
				if prev >= 'a' && prev <= 'z' {
					result = append(result, '_')
				} else if prev >= 'A' && prev <= 'Z' && i+1 < len(s) && s[i+1] >= 'a' && s[i+1] <= 'z' {
					result = append(result, '_')
				}
			}
			result = append(result, byte(c)+32)
		} else {
			result = append(result, byte(c))
		}
	}
	return string(result)
}
