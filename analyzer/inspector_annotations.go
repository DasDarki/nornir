package analyzer

import (
	"go/ast"
	"nornir/annotations"
	"strings"
)

type AnnotationsInspector struct{}

func (i *AnnotationsInspector) visit(ctx *InspectorContext, node ast.Node) bool {
	switch n := node.(type) {
	case *ast.Comment:
		annotation := annotations.Parse(strings.TrimSpace(n.Text))
		if annotation != nil {
		}
		break
	}
	return true
}

func (i *AnnotationsInspector) finish(a *Analyzer) error {
	return nil
}
