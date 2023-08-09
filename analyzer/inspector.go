package analyzer

import (
	"go/ast"
)

type InspectorContext struct {
	Analyzer *Analyzer
	Path     string
}

type inspector interface {
	visit(ctx *InspectorContext, node ast.Node) bool
	finish(a *Analyzer) error
}
