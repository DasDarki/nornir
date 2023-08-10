package generator

import (
	"go/ast"
	"nornir/analyzer"
	"nornir/log"
	"strconv"
	"strings"
)

func generateHandler(r *route, c *controller) ([]string, bool) {
	strconv := false
	url := r.Path
	if c != nil {
		url = strings.TrimPrefix(url, "/")

		if !strings.HasSuffix(c.Path, "/") {
			c.Path += "/"
		}

		url = c.Path + url
	}

	content := []string{
		"",
		"func " + cfg.Prefix + r.Funcname + "(c *gin.Context) {",
	}

	body := make([]string, 0)
	results := make([]string, 0)
	resultVars := make([]string, 0)
	var retExpr ast.Expr

	for _, ret := range r.Signature.Results.List {
		switch v := ret.Type.(type) {
		case *ast.Ident:
			switch v.Name {
			case "int":
				results = append(results, "int")
				resultVars = append(resultVars, "res")
			case "error":
				results = append(results, "error")
				resultVars = append(resultVars, "err")
			default:
				results = append(results, "any")
				resultVars = append(resultVars, "res")
				retExpr = v
			}
		default:
			results = append(results, "any")
			resultVars = append(resultVars, "res")
			retExpr = v
		}
	}

	ret := ""
	if len(resultVars) > 0 {
		ret = strings.Join(resultVars, ", ") + " := "
	}

	precode, callParams, usedStrconv := generatePreHandlerCode(r, url)
	body = addArray(body, precode...)
	body = addArray(body, ret+r.Funcname+"("+strings.Join(callParams, ", ")+")")

	if usedStrconv {
		strconv = true
	}

	switch len(results) {
	case 0:
		body = append(body, "c.JSON(http.StatusNoContent, nil)")
	case 1:
		switch results[0] {
		case "int":
			body = append(body, "c.JSON(res, nil)")
		case "any":
			body = append(body, "c.JSON(http.StatusOK, res)")
		case "error":
			body = append(body, "if err != nil {")
			body = append(body, "\tc.JSON(http.StatusInternalServerError, err)")
			body = append(body, "\treturn")
			body = append(body, "}")
			body = append(body, "c.JSON(http.StatusNoContent, nil)")
		}
	case 2:
		switch {
		case results[0] == "int" && results[1] == "error":
			body = append(body, "if err != nil {")
			body = append(body, "\tc.JSON(http.StatusInternalServerError, err)")
			body = append(body, "\treturn")
			body = append(body, "}")
			body = append(body, "c.JSON(res, nil)")
		case results[0] == "any" && results[1] == "error":
			body = append(body, "if err != nil {")
			body = append(body, "\tc.JSON(http.StatusInternalServerError, err)")
			body = append(body, "\treturn")
			body = append(body, "}")
			body = append(body, "c.JSON(http.StatusOK, res)")
		}
	default:
		panic("Too many return values")
	}

	if retExpr != nil && r.Usage != nil {
		tstype := generateDtoFromExpr(retExpr, r.Meta.Package, r.Meta.File)

		r.Usage.Response = &tstype
	}

	content = addArray(content, appendToStrings(body, "\t")...)
	return addArray(content, "}"), strconv
}

func generateDtoFromExpr(expr ast.Expr, currPackage string, file *ast.File) string {
	tstype := ""
	importMap := collectImportMap(file)
	importStruct := &analyzer.Struct{
		Imports: importMap,
	}

	switch v := expr.(type) {
	case *ast.Ident:
		possibleStruct := currPackage + "@" + v.Name
		if i, ok := structCache[possibleStruct]; ok {
			tstype = generateTypeScriptObjectBody(i.Node, &i, currPackage)
			break
		}

		tstype = translateGoPrimitivesToTypeScript(v.Name)
	case *ast.StarExpr:
		tstype = generateTypeScriptType(v.X, importStruct, currPackage) + " | null"
	case *ast.ArrayType:
		tstype = generateTypeScriptType(v.Elt, importStruct, currPackage) + "[]"
	case *ast.MapType:
		tstype = "Map<" + generateTypeScriptType(v.Key, importStruct, currPackage) + ", " + generateTypeScriptType(v.Value, importStruct, currPackage) + ">"
	case *ast.SelectorExpr:
		tstype = generateDtoFromSelector(currPackage, file, v.X.(*ast.Ident).Name, v.Sel.Name)
	case *ast.InterfaceType:
		tstype = "any"
	case *ast.StructType:
		tstype = "(" + generateTypeScriptObjectBody(v, importStruct, currPackage) + ")"
	default:
		log.Debugf("Unknown type: %T", expr)
		tstype = "unknown"
	}

	return tstype
}

func generateDtoFromSelector(currPackage string, file *ast.File, pkgName string, typeName string) string {
	importMap := collectImportMap(file)

	var currStruct *analyzer.Struct = nil

	importedPackage, ok := importMap[pkgName]
	if !ok {
		log.Debugf("gingen: Failed to find import for %s", pkgName)
	} else {
		importedStruct, ok := structCache[importedPackage+"@"+typeName]
		if !ok {
			log.Debugf("gingen: Failed to find struct for %s", importedPackage+"@"+typeName)
		} else {
			currStruct = &importedStruct
		}
	}

	if currStruct != nil {
		return generateTypeScriptObjectBody(currStruct.Node, currStruct, currPackage)
	}

	return "any"
}

func generatePreHandlerCode(r *route, fullurl string) ([]string, []string, bool) {
	body := []string{}
	params := []string{}
	usedStrconv := false

	r.Usage = &usage{
		FullPath: fullurl,
		Method:   r.Method,
		Name:     r.Funcname,
		Params:   []string{},
		Header:   []string{},
		Query:    map[string]string{},
		Body:     nil,
	}

	for _, param := range r.Signature.Params.List {
		name := param.Names[0].Name
		typeName := simplifyType(param.Type)

		if strings.Contains(typeName, "*gin.Context") {
			params = append(params, "c")
			continue
		}

		if strings.Contains(fullurl, ":"+name) {
			body = append(body, "p_"+name+" := c.Param(\""+name+"\")")
			body = append(body, "if p_"+name+" == \"\" {")
			body = append(body, "\tc.JSON(http.StatusBadRequest, \"Missing parameter "+name+"\")")
			body = append(body, "\treturn")
			body = append(body, "}")
			body = append(body, "")

			params = append(params, "p_"+name)

			r.Usage.Params = append(r.Usage.Params, name)
			continue
		}

		if strings.HasPrefix(name, "query") {
			body = append(body, "q_"+name+" := c.Query(\""+name[5:]+"\")")
			body = append(body, "if q_"+name+" == \"\" {")
			body = append(body, "\tc.JSON(http.StatusBadRequest, \"Missing parameter "+name[5:]+"\")")
			body = append(body, "\treturn")
			body = append(body, "}")
			body = append(body, "")

			convCode, jsType := getConvertCodeForQuery(typeName, "q_"+name, "qc_"+name)
			body = append(body, convCode...)
			usedStrconv = true

			params = append(params, "qc_"+name)

			r.Usage.Query[name[5:]] = jsType
			continue
		} else if strings.HasPrefix(name, "header") {
			body = append(body, "h_"+name+" := c.GetHeader(\""+name[6:]+"\")")
			body = append(body, "if h_"+name+" == \"\" {")
			body = append(body, "\tc.JSON(http.StatusBadRequest, \"Missing parameter "+name[6:]+"\")")
			body = append(body, "\treturn")
			body = append(body, "}")
			body = append(body, "")

			params = append(params, "h_"+name)

			r.Usage.Header = append(r.Usage.Header, name[6:])
			continue
		} else if strings.HasPrefix(name, "body") {
			isPtr := false
			if strings.HasPrefix(typeName, "*") {
				isPtr = true
				typeName = typeName[1:]
			}

			body = append(body, "b_"+name+" := "+typeName+"{}")
			body = append(body, "if err := c.BindJSON(&b_"+name+"); err != nil {")
			body = append(body, "\tc.JSON(http.StatusBadRequest, \"Invalid parameter "+name+"\")")
			body = append(body, "\treturn")
			body = append(body, "}")
			body = append(body, "")

			prepend := ""
			if isPtr {
				prepend = "&"
			}

			params = append(params, prepend+"b_"+name)

			var tstype string = ""

			if strings.Contains(typeName, ".") {
				parts := strings.Split(typeName, ".")
				pkgName := parts[0]
				typeName := parts[1]
				tstype = generateDtoFromSelector(r.Meta.Package, r.Meta.File, pkgName, typeName)
			} else {
				possibleStruct := r.Meta.Package + "@" + typeName
				if i, ok := structCache[possibleStruct]; ok {
					tstype = generateTypeScriptObjectBody(i.Node, &i, r.Meta.Package)
				} else {
					tstype = translateGoPrimitivesToTypeScript(typeName)
				}
			}

			r.Usage.Body = &tstype
			continue
		} else if isQueryDefaultDataSource(r) {
			body = append(body, "q_"+name+" := c.Query(\""+name+"\")")
			body = append(body, "if q_"+name+" == \"\" {")
			body = append(body, "\tc.JSON(http.StatusBadRequest, \"Missing parameter "+name+"\")")
			body = append(body, "\treturn")
			body = append(body, "}")
			body = append(body, "")

			convCode, jsType := getConvertCodeForQuery(typeName, "q_"+name, "qc_"+name)
			body = append(body, convCode...)
			usedStrconv = true

			params = append(params, "qc_"+name)

			r.Usage.Query[name[5:]] = jsType
			continue
		} else if isBodyDefaultDataSource(r) {
			isPtr := false
			if strings.HasPrefix(typeName, "*") {
				isPtr = true
				typeName = typeName[1:]
			}

			body = append(body, "b_"+name+" := "+typeName+"{}")
			body = append(body, "if err := c.BindJSON(&b_"+name+"); err != nil {")
			body = append(body, "\tc.JSON(http.StatusBadRequest, \"Invalid parameter "+name+"\")")
			body = append(body, "\treturn")
			body = append(body, "}")
			body = append(body, "")

			prepend := ""
			if isPtr {
				prepend = "&"
			}

			params = append(params, prepend+"b_"+name)

			var tstype string = ""

			if strings.Contains(typeName, ".") {
				parts := strings.Split(typeName, ".")
				pkgName := parts[0]
				typeName := parts[1]
				tstype = generateDtoFromSelector(r.Meta.Package, r.Meta.File, pkgName, typeName)
			} else {
				possibleStruct := r.Meta.Package + "@" + typeName
				if i, ok := structCache[possibleStruct]; ok {
					tstype = generateTypeScriptObjectBody(i.Node, &i, r.Meta.Package)
				} else {
					tstype = translateGoPrimitivesToTypeScript(typeName)
				}
			}

			r.Usage.Body = &tstype
			continue
		}
	}

	return body, params, usedStrconv
}

func isQueryDefaultDataSource(r *route) bool {
	if r.Method == "GET" || r.Method == "DELETE" {
		return true
	}

	return false
}

func isBodyDefaultDataSource(r *route) bool {
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
		return true
	}

	return false
}

func simplifyType(t ast.Expr) string {
	switch t := t.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + simplifyType(t.X)
	case *ast.ArrayType:
		return "[]" + simplifyType(t.Elt)
	case *ast.MapType:
		return "map[" + simplifyType(t.Key) + "]" + simplifyType(t.Value)
	case *ast.SelectorExpr:
		return t.X.(*ast.Ident).Name + "." + t.Sel.Name
	default:
		panic("Unknown type")
	}
}

func getConvertCodeForQuery(t string, inParam string, outParam string) ([]string, string) {
	if t == "string" {
		return []string{outParam + " := " + inParam}, "string"
	}

	if t == "int" {
		return []string{outParam + ", err := strconv.Atoi(" + inParam + ")", "if err != nil {", "\tc.JSON(http.StatusBadRequest, \"Invalid parameter " + inParam + "\")", "\treturn", "}"}, "number"
	}

	if t == "int64" {
		return []string{outParam + ", err := strconv.ParseInt(" + inParam + ", 10, 64)", "if err != nil {", "\tc.JSON(http.StatusBadRequest, \"Invalid parameter " + inParam + "\")", "\treturn", "}"}, "number"
	}

	if t == "float64" {
		return []string{outParam + ", err := strconv.ParseFloat(" + inParam + ", 64)", "if err != nil {", "\tc.JSON(http.StatusBadRequest, \"Invalid parameter " + inParam + "\")", "\treturn", "}"}, "number"
	}

	if t == "bool" {
		return []string{outParam + ", err := strconv.ParseBool(" + inParam + ")", "if err != nil {", "\tc.JSON(http.StatusBadRequest, \"Invalid parameter " + inParam + "\")", "\treturn", "}"}, "boolean"
	}

	panic("Unsupported type " + t + " for query parameter")
}

func collectImportMap(file *ast.File) map[string]string {
	importMap := make(map[string]string)

	for _, imp := range file.Imports {
		path, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			panic(err)
		}

		if imp.Name != nil {
			importMap[imp.Name.Name] = path
		} else {
			parts := strings.Split(path, "/")
			importMap[parts[len(parts)-1]] = path
		}
	}

	return importMap
}
