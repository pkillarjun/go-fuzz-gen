package main

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
)

func harness_function(files_ast []*ast.File) []*ast.FuncDecl {
	var harness_funcs []*ast.FuncDecl

	for _, file_ast := range files_ast {
		for _, decl := range file_ast.Decls {

			/* check if 'decl' is a 'ast.FuncDecl'. */
			func_decl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			/*
			 * List all the functions that accept parameters as:
			 * string, []byte, io.Reader
			 */
			for _, field := range func_decl.Type.Params.List {
				switch field_type := field.Type.(type) {

				case *ast.Ident:
					if field_type.Name == "string" {
						harness_funcs = append(harness_funcs, func_decl)
						break
					}
				case *ast.ArrayType:
					array_type, ok := field_type.Elt.(*ast.Ident)
					if !ok {
						continue
					}

					if array_type.Name == "byte" {
						harness_funcs = append(harness_funcs, func_decl)
						break
					}
				case *ast.SelectorExpr:
					x_ident, ok := field_type.X.(*ast.Ident)
					if !ok {
						continue
					}

					if (x_ident.Name == "io") && (field_type.Sel.Name == "Reader") {
						harness_funcs = append(harness_funcs, func_decl)
						break
					}
				}
			}
		}
	}

	return harness_funcs
}

func node_string(node ast.Node) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), node)
	if err != nil {
		return ""
	}
	return buf.String()
}

func ast_analysis(codeast []*ast.File) {

	func_map := make(map[string][]*ast.FuncDecl)

	harness_funcs := harness_function(codeast)
	if len(harness_funcs) == 0 {
		return
	}

	for _, harness_func := range harness_funcs {
		func_name := harness_func.Name.Name
		func_map[func_name] = append(func_map[func_name], harness_func)
	}

	for func_name, func_decls := range func_map {
		var functions string

		/*
		 * Get all the implementations of a function with one name.
		 */
		for _, func_decl := range func_decls {
			functions += node_string(func_decl) + "\n"
		}
		run_gpt(func_name, functions)
	}
}
