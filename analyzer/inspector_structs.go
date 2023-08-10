package analyzer

import (
	"go/ast"
	"nornir/log"
	"strings"
)

type StructInspector struct{}

func (i *StructInspector) visit(ctx *InspectorContext, node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if ok {
		for _, spec := range decl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			structPath := ctx.Path + "@" + typeSpec.Name.Name
			ctx.Analyzer.Structs[structPath] = Struct{
				Node:    structType,
				Name:    typeSpec.Name.Name,
				Imports: make(map[string]string),
			}
		}
	}

	importDecl, ok := node.(*ast.ImportSpec)
	if ok {
		if ctx.Analyzer.StructImportMap[ctx.Path] == nil {
			ctx.Analyzer.StructImportMap[ctx.Path] = []Import{}
		}

		path := importDecl.Path.Value[1 : len(importDecl.Path.Value)-1]
		parts := strings.Split(path, "/")
		name := parts[len(parts)-1]

		if importDecl.Name != nil {
			name = importDecl.Name.Name
		}

		ctx.Analyzer.StructImportMap[ctx.Path] = append(ctx.Analyzer.StructImportMap[ctx.Path], Import{
			Path: path,
			Name: name,
		})
	}

	return true
}

func (i *StructInspector) finish(a *Analyzer) error {
	for structPath, cachedStruct := range a.Structs {
		log.Debugf("Found struct at %s (%s)", structPath, cachedStruct.Name)

		parts := strings.Split(structPath, "@")
		pkgPath := parts[0]

		for pkgPath2, imports := range a.StructImportMap {
			if pkgPath2 == pkgPath {
				for _, importDecl := range imports {
					cachedStruct.Imports[importDecl.Name] = importDecl.Path
				}
			}
		}
	}

	for structPath, imports := range a.StructImportMap {
		log.Debugf("Found imports for %s: %v", structPath, imports)
	}

	return nil
}
