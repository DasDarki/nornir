package analyzer

import "go/ast"

type AnnotationsInspector struct{}

func (i *AnnotationsInspector) visit(ctx *InspectorContext, node ast.Node) bool {
	return true
}

func (i *AnnotationsInspector) finish(a *Analyzer) error {
	return nil
}
