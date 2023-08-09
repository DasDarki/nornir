package analyzer

import (
	"go/ast"
)

type InspectorContext struct {
	Analyzer *Analyzer
	Path     string
	FileId   string
	Filename string
	Data     map[string]interface{}
}

type inspector interface {
	visit(ctx *InspectorContext, node ast.Node) bool
	finish(a *Analyzer) error
}
