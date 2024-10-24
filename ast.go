package main

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
)

func hasTargetParam(params []*ast.Field) bool {
	for _, field := range params {
		switch fieldType := field.Type.(type) {
		case *ast.Ident:
			if fieldType.Name == "string" {
				return true
			}
		case *ast.ArrayType:
			if eltIdent, ok := fieldType.Elt.(*ast.Ident); ok && eltIdent.Name == "byte" {
				return true
			}
		case *ast.SelectorExpr:
			if xIdent, ok := fieldType.X.(*ast.Ident); ok && xIdent.Name == "io" && fieldType.Sel.Name == "Reader" {
				return true
			}
		}
	}
	return false
}

func getHarnessFuncs(filesAST []*ast.File) []*ast.FuncDecl {
	var harnessFuncs []*ast.FuncDecl

	for _, fileAST := range filesAST {
		for _, decl := range fileAST.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			if hasTargetParam(funcDecl.Type.Params.List) {
				harnessFuncs = append(harnessFuncs, funcDecl)
			}
		}
	}
	return harnessFuncs
}

func processAST(filesAST []*ast.File) {
	harnessFuncs := getHarnessFuncs(filesAST)
	if len(harnessFuncs) == 0 {
		return
	}

	for _, funcDecl := range harnessFuncs {
		var buf bytes.Buffer
		printer.Fprint(&buf, token.NewFileSet(), funcDecl)
		funcName := funcDecl.Name.Name
		processGPT(funcName, buf.String())
	}
}
