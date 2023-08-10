package generator

import (
	"encoding/json"
	"go/ast"
	"nornir/analyzer"
	"nornir/log"
	"strings"
)

var dtoRegistry = make(map[string]string)
var structCache = make(map[string]analyzer.Struct)

func GenerateDTOs(a *analyzer.Analyzer) {
	for structPath, strct := range a.Structs {
		structCache[structPath] = strct
	}

	for structPath, strct := range a.Structs {
		dtoRegistry[structPath] = generateTypeScriptObjectBody(strct.Node, &strct, strings.Split(structPath, "@")[0])
	}

	log.Debugf("Generated %d DTOs", len(dtoRegistry))
	json, err := json.Marshal(dtoRegistry)
	if err != nil {
		log.Debugf("Failed to marshal DTO registry: %s", err)
		return
	}

	log.Debugf("%s", string(json))
}

func generateTypeScriptType(expr ast.Expr, currStruct *analyzer.Struct, currPackage string) string {
	switch v := expr.(type) {
	case *ast.Ident:
		possibleStruct := currPackage + "@" + v.Name
		if i, ok := structCache[possibleStruct]; ok {
			return generateTypeScriptObjectBody(i.Node, &i, currPackage)
		}

		return translateGoPrimitivesToTypeScript(v.Name)
	case *ast.StarExpr:
		return generateTypeScriptType(v.X, currStruct, currPackage) + " | null"
	case *ast.ArrayType:
		return generateTypeScriptType(v.Elt, currStruct, currPackage) + "[]"
	case *ast.MapType:
		return "Map<" + generateTypeScriptType(v.Key, currStruct, currPackage) + ", " + generateTypeScriptType(v.Value, currStruct, currPackage) + ">"
	case *ast.SelectorExpr:
		pkgName := v.X.(*ast.Ident).Name
		typeName := v.Sel.Name

		log.Debugf("pkgName: %s", pkgName)
		log.Debugf("typeName: %s", typeName)

		importedPackage, ok := currStruct.Imports[pkgName]
		if !ok {
			log.Debugf("Failed to find import for %s", pkgName)
			return "any"
		}

		importedStruct, ok := structCache[importedPackage+"@"+typeName]
		if !ok {
			log.Debugf("Failed to find struct for %s", importedPackage+"@"+typeName)
			return "any"
		}

		return generateTypeScriptObjectBody(importedStruct.Node, &importedStruct, currPackage)
	case *ast.InterfaceType:
		return "any"
	case *ast.StructType:
		return "(" + generateTypeScriptObjectBody(v, currStruct, currPackage) + ")"
	default:
		log.Debugf("Unknown type: %T", expr)
		return "unknown"
	}
}

func generateTypeScriptObjectBody(strct *ast.StructType, currStruct *analyzer.Struct, currPackage string) string {
	lines := make([]string, 0)

	for _, field := range strct.Fields.List {
		name := getFieldName(field)
		if name == "" {
			continue
		}

		line := name + ": " + generateTypeScriptType(field.Type, currStruct, currPackage)
		lines = append(lines, line)
	}

	return "{" + strings.Join(lines, ", ") + "}"
}

func getFieldName(field *ast.Field) string {
	name := ""
	if field.Tag != nil {
		if !strings.Contains(field.Tag.Value, "json") {
			return ""
		}

		tag := strings.Trim(field.Tag.Value, "`")
		tag = strings.TrimPrefix(tag, "json:\"")
		tag = strings.TrimSuffix(tag, "\"")
		tagParts := strings.Split(tag, ",")
		name = tagParts[0]
	}

	if name == "" {
		name = makeNameLowerCamelCase(field.Names[0].Name)
	}

	return name
}

func translateGoPrimitivesToTypeScript(primitive string) string {
	switch primitive {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint32", "uint64", "uintptr", "float32", "float64", "complex64", "complex128", "byte", "rune":
		return "number"
	case "string":
		return "string"
	case "bool":
		return "boolean"
	}

	return primitive
}

func makeNameLowerCamelCase(name string) string {
	return strings.ToLower(name[0:1]) + name[1:]
}
