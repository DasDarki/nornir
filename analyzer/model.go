package analyzer

import "go/ast"

type Struct struct {
	Node *ast.StructType
	Name string
}

type Import struct {
	Path string
	Name string
}
