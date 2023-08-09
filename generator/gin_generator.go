package generator

import (
	"encoding/json"
	"go/ast"
	"nornir/analyzer"
	"nornir/annotations"
	"nornir/log"
	"os"
	"path/filepath"
)

type app struct {
	controller
	ControllerList []*controller `json:"controllers"`
}

type controller struct {
	Dir    string   `json:"dir"`
	Pack   string   `json:"package"`
	Path   string   `json:"path"`
	Routes []*route `json:"routes"`
}

type route struct {
	Method    string        `json:"method"`
	Path      string        `json:"path"`
	Signature *ast.FuncType `json:"-"`
	Funcname  string        `json:"funcname"`
}

func GenerateGinCode(a *analyzer.Analyzer) {
	app := structureApp(a)

	yamlText, err := json.Marshal(app)
	if err != nil {
		panic(err)
	}

	log.Debugf("Generating gin code")
	log.Debugf("JSON: %s", string(yamlText))
}

func structureApp(a *analyzer.Analyzer) *app {
	currDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	app := &app{
		controller: controller{
			Dir:    currDir,
			Pack:   a.ModName,
			Path:   "",
			Routes: []*route{},
		},
		ControllerList: []*controller{},
	}

	createdControllers := make(map[string]*controller)
	finishedFiles := make(map[string]bool)

	for _, file := range a.AnnotatedFiles {
		for _, fileAnnotation := range file.Annotations {
			if controllerAnnotation, ok := fileAnnotation.(*annotations.ControllerAnnotation); ok {
				controllerKey := file.Meta.Package + "@" + controllerAnnotation.Path
				if _, ok := createdControllers[controllerKey]; !ok {
					controller := &controller{
						Dir:    getDirectoryOfFile(file.Meta.Filename),
						Pack:   file.Meta.Package,
						Path:   removeQuotes(controllerAnnotation.Path),
						Routes: []*route{},
					}
					app.ControllerList = append(app.ControllerList, controller)
					createdControllers[controllerKey] = controller
				}

				controller := createdControllers[controllerKey]
				for _, function := range a.AnnotatedFuncs {
					if function.Meta.FileId == file.Meta.FileId {
						addRoutesToController(controller, function)
						finishedFiles[file.Meta.FileId] = true
					}
				}
			}
		}
	}

	for _, function := range a.AnnotatedFuncs {
		if _, ok := finishedFiles[function.Meta.FileId]; !ok {
			addRoutesToController(&app.controller, function)
		}
	}

	return app
}

func addRoutesToController(controller *controller, function *annotations.AnnotatedFunc) {
	for _, funcAnnotation := range function.Annotations {
		if requestMappingAnnotation, ok := funcAnnotation.(*annotations.RequestMappingAnnotation); ok {
			controller.Routes = append(controller.Routes, &route{
				Method:    requestMappingAnnotation.Method,
				Path:      removeQuotes(requestMappingAnnotation.Path),
				Signature: function.Func.Type,
				Funcname:  function.Func.Name.Name,
			})
		} else if getAnnotation, ok := funcAnnotation.(*annotations.GetAnnotation); ok {
			controller.Routes = append(controller.Routes, &route{
				Method:    "GET",
				Path:      removeQuotes(getAnnotation.Path),
				Signature: function.Func.Type,
				Funcname:  function.Func.Name.Name,
			})
		} else if postAnnotation, ok := funcAnnotation.(*annotations.PostAnnotation); ok {
			controller.Routes = append(controller.Routes, &route{
				Method:    "POST",
				Path:      removeQuotes(postAnnotation.Path),
				Signature: function.Func.Type,
				Funcname:  function.Func.Name.Name,
			})
		} else if putAnnotation, ok := funcAnnotation.(*annotations.PutAnnotation); ok {
			controller.Routes = append(controller.Routes, &route{
				Method:    "PUT",
				Path:      removeQuotes(putAnnotation.Path),
				Signature: function.Func.Type,
				Funcname:  function.Func.Name.Name,
			})
		} else if deleteAnnotation, ok := funcAnnotation.(*annotations.DeleteAnnotation); ok {
			controller.Routes = append(controller.Routes, &route{
				Method:    "DELETE",
				Path:      removeQuotes(deleteAnnotation.Path),
				Signature: function.Func.Type,
				Funcname:  function.Func.Name.Name,
			})
		} else if patchAnnotation, ok := funcAnnotation.(*annotations.PatchAnnotation); ok {
			controller.Routes = append(controller.Routes, &route{
				Method:    "PATCH",
				Path:      removeQuotes(patchAnnotation.Path),
				Signature: function.Func.Type,
				Funcname:  function.Func.Name.Name,
			})
		} else if headAnnotation, ok := funcAnnotation.(*annotations.HeadAnnotation); ok {
			controller.Routes = append(controller.Routes, &route{
				Method:    "HEAD",
				Path:      removeQuotes(headAnnotation.Path),
				Signature: function.Func.Type,
				Funcname:  function.Func.Name.Name,
			})
		} else if optionsAnnotation, ok := funcAnnotation.(*annotations.OptionsAnnotation); ok {
			controller.Routes = append(controller.Routes, &route{
				Method:    "OPTIONS",
				Path:      removeQuotes(optionsAnnotation.Path),
				Signature: function.Func.Type,
				Funcname:  function.Func.Name.Name,
			})
		}
	}
}

func removeQuotes(s string) string {
	if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}

	return s
}

func getDirectoryOfFile(filename string) string {
	return filepath.Dir(filename)
}
