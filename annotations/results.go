package annotations

import "go/ast"

type Meta struct {
	Package  string
	FileId   string
	Filename string
}

type AnnotatedFile struct {
	Meta        Meta
	Annotations []interface{}
	File        *ast.File
}

type AnnotatedFunc struct {
	Meta        Meta
	Annotations []interface{}
	Func        *ast.FuncDecl
}

type AnnotatedVar struct {
	Meta        Meta
	Annotations []interface{}
	Var         *ast.ValueSpec
}
