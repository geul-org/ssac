package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/ssac/artifacts/internal/generator"
	"github.com/park-jun-woo/ssac/artifacts/internal/parser"
	"github.com/park-jun-woo/ssac/artifacts/internal/validator"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: ssac <command>")
		fmt.Fprintln(os.Stderr, "commands: parse, validate, gen")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "parse":
		runParse()
	case "validate":
		runValidate()
	case "gen":
		runGen()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runValidate() {
	dir := "specs/backend/service"
	if len(os.Args) > 2 {
		dir = os.Args[2]
	}

	// dir/service/ 가 있으면 프로젝트 루트, 없으면 service 디렉토리 직접 지정
	serviceDir := filepath.Join(dir, "service")
	projectRoot := dir
	if _, err := os.Stat(serviceDir); os.IsNotExist(err) {
		serviceDir = dir
		projectRoot = ""
	}

	funcs, err := parser.ParseDir(serviceDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	// 프로젝트 루트가 있고 외부 SSOT가 있으면 심볼 테이블 교차 검증
	if projectRoot != "" {
		st, stErr := validator.LoadSymbolTable(projectRoot)
		if stErr == nil {
			errs := validator.ValidateWithSymbols(funcs, st)
			if len(errs) == 0 {
				fmt.Println("validation passed (with symbol table)")
				return
			}
			for _, e := range errs {
				fmt.Fprintf(os.Stderr, "ERROR: %s\n", e)
			}
			os.Exit(1)
		}
	}

	// 외부 SSOT 없으면 내부 검증만
	errs := validator.Validate(funcs)
	if len(errs) == 0 {
		fmt.Println("validation passed")
		return
	}

	for _, e := range errs {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", e)
	}
	os.Exit(1)
}

func runGen() {
	inDir := "specs/backend/service"
	outDir := "artifacts/backend/internal/service"

	funcs, err := parser.ParseDir(inDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	// validate before generate
	errs := validator.Validate(funcs)
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", e)
		}
		fmt.Fprintln(os.Stderr, "validation failed, code generation aborted")
		os.Exit(1)
	}

	if err := generator.Generate(funcs, outDir); err != nil {
		fmt.Fprintf(os.Stderr, "generate error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("generated %d files in %s\n", len(funcs), outDir)
}

func runParse() {
	dir := "specs/backend/service"
	if len(os.Args) > 2 {
		dir = os.Args[2]
	}

	funcs, err := parser.ParseDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	for _, f := range funcs {
		fmt.Printf("=== %s (%s) ===\n", f.Name, f.FileName)
		for i, s := range f.Sequences {
			fmt.Printf("  [%d] %s", i, s.Type)
			if s.Model != "" {
				fmt.Printf(" | model=%s", s.Model)
			}
			if s.Result != nil {
				fmt.Printf(" | result=%s %s", s.Result.Var, s.Result.Type)
			}
			if s.Message != "" {
				fmt.Printf(" | message=%q", s.Message)
			}
			fmt.Println()
		}
	}
}
