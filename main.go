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
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorBlue    = "\033[34m"
	colorCyan    = "\033[36m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorMagenta = "\033[35m"
)

func processDirectory(directory string) {
	var goFiles []string
	var filesAST []*ast.File

	entries, err := os.ReadDir(directory)
	if err != nil {
		fmt.Printf("%sError reading directory '%s': %v%s\n", colorRed, directory, err, colorReset)
		return
	}

	for _, entry := range entries {
		entryPath := filepath.Join(directory, entry.Name())

		if entry.IsDir() {
			processDirectory(entryPath)
			continue
		}

		if strings.HasSuffix(entry.Name(), ".go") && !strings.HasSuffix(entry.Name(), "_test.go") {
			goFiles = append(goFiles, entryPath)
		}
	}

	if len(goFiles) == 0 {
		return
	}

	fmt.Printf("%sProcessing directory: %s%s\n", colorBlue, directory, colorReset)

	fset := token.NewFileSet()
	for _, filePath := range goFiles {
		fileAST, err := parser.ParseFile(fset, filePath, nil, 0)
		if err != nil {
			fmt.Printf("%sError parsing file '%s': %v%s\n", colorRed, filePath, err, colorReset)
			continue
		}
		filesAST = append(filesAST, fileAST)
	}

	processAST(filesAST)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("%sUsage: %s <directory>%s\n", colorRed, os.Args[0], colorReset)
		return
	}

	processDirectory(os.Args[1])
}
