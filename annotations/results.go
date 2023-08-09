package annotations

import "go/ast"

type AnnotatedPackage struct {
	Annotations []interface{}
	Package     *ast.File
}

type AnnotatedFunc struct {
	Annotations []interface{}
	Func        *ast.FuncDecl
}

type AnnotatedVar struct {
	Annotations []interface{}
	Var         *ast.ValueSpec
}
