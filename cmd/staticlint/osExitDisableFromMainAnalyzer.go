package main

import (
	"errors"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// ExitDisableAnalyzer analyzer
var ExitDisableAnalyzer = &analysis.Analyzer{
	Name: "osExitDisableFromMain",
	Doc:  "check for disable the use of a direct call to os.Exit in the main package main function",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name == "main" {
			for _, decl := range file.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok {
					if fn.Name.Name == "main" {
						for _, stmt := range fn.Body.List {
							if call, ok := stmt.(*ast.ExprStmt); ok {
								if sel, ok := call.X.(*ast.SelectorExpr); ok {
									if sel.Sel.Name == "Exit" {
										return nil, errors.New("direct call to os.Exit in main package main function detected")
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return nil, nil
}
