package analyzer

import "go/ast"

type Struct struct {
	Node    *ast.StructType   `json:"-"`
	Name    string            `json:"name"`
	Imports map[string]string `json:"imports"`
}

type Import struct {
	Path string
	Name string
}
