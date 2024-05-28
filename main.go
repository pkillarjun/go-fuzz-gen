package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

const (
	red_color     = "\033[31m"
	blue_color    = "\033[34m"
	cyan_color    = "\033[36m"
	green_color   = "\033[32m"
	reset_color   = "\033[0m"
	yellow_color  = "\033[33m"
	magenta_color = "\033[35m"
)

/*
 * Check all source files and *_test.go files for fuzz tests.
 * Get AST of files.
 */
func process_files(files []string) {

	var files_ast []*ast.File
	fset := token.NewFileSet()

	for _, file := range files {
		file_ast, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			fmt.Printf(red_color+"error: %s: %v"+reset_color, file, err)
			continue
		}

		files_ast = append(files_ast, file_ast)
	}

	ast_analysis(files_ast)
}

func process_dir(path string) {
	var files []string

	dir_entry, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf(red_color+"Error reading directory: %v\n"+reset_color, err)
		return
	}

	for _, entry := range dir_entry {
		entry_path := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			process_dir(entry_path)
			continue
		}

		if strings.HasSuffix(entry.Name(), ".go") {
			files = append(files, entry_path)
		}
	}

	if len(files) > 0 {
		fmt.Printf(blue_color+"\nProcessing directory: %s\n"+reset_color, path)
		process_files(files)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf(red_color + "Provide input as <dir> \n" + reset_color)
		return
	}

	process_dir(os.Args[1])
}
