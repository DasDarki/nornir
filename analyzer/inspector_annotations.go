package analyzer

import (
	"go/ast"
	"log"
	"nornir/annotations"
)

type AnnotationsInspector struct{}

func (i *AnnotationsInspector) visit(ctx *InspectorContext, node ast.Node) bool {
	switch n := node.(type) {
	case *ast.File:
		if n.Doc != nil {
			annotationList := make([]interface{}, 0)

			for _, c := range n.Doc.List {
				annotation := annotations.Parse(c.Text)
				if annotation != nil {
					if annotation.Name == "Controller" {
						controllerAnnotation := &annotations.ControllerAnnotation{}
						err := annotations.Decode(annotation, controllerAnnotation, annotations.UsageKindPackage)
						if err != nil {
							log.Println("WARNING: Failed to decode Controller annotation: " + err.Error())
							continue
						}

						annotationList = append(annotationList, controllerAnnotation)
					} else {
						log.Println("WARNING: Unknown package annotation: " + annotation.Name)
						continue
					}
				}
			}

			if len(annotationList) > 0 {
				ctx.Analyzer.AnnotatedFiles = append(ctx.Analyzer.AnnotatedFiles, &annotations.AnnotatedFile{
					Meta:        annotations.Meta{Package: ctx.Path, FileId: ctx.FileId, Filename: ctx.Filename},
					Annotations: annotationList,
					File:        n,
				})
			}
		}

		for _, decl := range n.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				if funcDecl.Doc != nil {
					annotationList := make([]interface{}, 0)

					for _, c := range funcDecl.Doc.List {
						annotation := annotations.Parse(c.Text)
						if annotation != nil {
							if annotation.Name == "RequestMapping" {
								requestMappingAnnotation := &annotations.RequestMappingAnnotation{}
								err := annotations.Decode(annotation, requestMappingAnnotation, annotations.UsageKindFunction)
								if err != nil {
									log.Println("WARNING: Failed to decode RequestMapping annotation: " + err.Error())
									continue
								}

								annotationList = append(annotationList, requestMappingAnnotation)
							} else if annotation.Name == "Get" {
								getAnnotation := &annotations.GetAnnotation{}
								err := annotations.Decode(annotation, getAnnotation, annotations.UsageKindFunction)
								if err != nil {
									log.Println("WARNING: Failed to decode Get annotation: " + err.Error())
									continue
								}

								annotationList = append(annotationList, getAnnotation)
							} else if annotation.Name == "Post" {
								postAnnotation := &annotations.PostAnnotation{}
								err := annotations.Decode(annotation, postAnnotation, annotations.UsageKindFunction)
								if err != nil {
									log.Println("WARNING: Failed to decode Post annotation: " + err.Error())
									continue
								}

								annotationList = append(annotationList, postAnnotation)
							} else if annotation.Name == "Put" {
								putAnnotation := &annotations.PutAnnotation{}
								err := annotations.Decode(annotation, putAnnotation, annotations.UsageKindFunction)
								if err != nil {
									log.Println("WARNING: Failed to decode Put annotation: " + err.Error())
									continue
								}

								annotationList = append(annotationList, putAnnotation)
							} else if annotation.Name == "Delete" {
								deleteAnnotation := &annotations.DeleteAnnotation{}
								err := annotations.Decode(annotation, deleteAnnotation, annotations.UsageKindFunction)
								if err != nil {
									log.Println("WARNING: Failed to decode Delete annotation: " + err.Error())
									continue
								}

								annotationList = append(annotationList, deleteAnnotation)
							} else if annotation.Name == "Patch" {
								patchAnnotation := &annotations.PatchAnnotation{}
								err := annotations.Decode(annotation, patchAnnotation, annotations.UsageKindFunction)
								if err != nil {
									log.Println("WARNING: Failed to decode Patch annotation: " + err.Error())
									continue
								}

								annotationList = append(annotationList, patchAnnotation)
							} else if annotation.Name == "Head" {
								headAnnotation := &annotations.HeadAnnotation{}
								err := annotations.Decode(annotation, headAnnotation, annotations.UsageKindFunction)
								if err != nil {
									log.Println("WARNING: Failed to decode Head annotation: " + err.Error())
									continue
								}

								annotationList = append(annotationList, headAnnotation)
							} else if annotation.Name == "Options" {
								optionsAnnotation := &annotations.OptionsAnnotation{}
								err := annotations.Decode(annotation, optionsAnnotation, annotations.UsageKindFunction)
								if err != nil {
									log.Println("WARNING: Failed to decode Options annotation: " + err.Error())
									continue
								}

								annotationList = append(annotationList, optionsAnnotation)
							} else {
								log.Println("WARNING: Unknown func annotation: " + annotation.Name)
								continue
							}
						}
					}

					if len(annotationList) > 0 {
						ctx.Analyzer.AnnotatedFuncs = append(ctx.Analyzer.AnnotatedFuncs, &annotations.AnnotatedFunc{
							Meta:        annotations.Meta{Package: ctx.Path, FileId: ctx.FileId, Filename: ctx.Filename, File: n},
							Annotations: annotationList,
							Func:        funcDecl,
						})
					}
				}
			}
		}
	case *ast.ValueSpec:
		if n.Doc != nil {
			annotationList := make([]interface{}, 0)

			for _, c := range n.Doc.List {
				annotation := annotations.Parse(c.Text)
				if annotation != nil {
					if annotation.Name == "InjectNornir" {
						injectAnnotation := &annotations.InjectNornirAnnotation{}
						err := annotations.Decode(annotation, injectAnnotation, annotations.UsageKindVar)
						if err != nil {
							log.Println("WARNING: Failed to decode InjectNornir annotation: " + err.Error())
							continue
						}

						annotationList = append(annotationList, injectAnnotation)
					} else {
						log.Println("WARNING: Unknown var annotation: " + annotation.Name)
						continue
					}
				}
			}

			if len(annotationList) > 0 {
				ctx.Analyzer.AnnotatedVars = append(ctx.Analyzer.AnnotatedVars, &annotations.AnnotatedVar{
					Meta:        annotations.Meta{Package: ctx.Path, FileId: ctx.FileId, Filename: ctx.Filename},
					Annotations: annotationList,
					Var:         n,
				})
			}
		}
	}
	return true
}

func (i *AnnotationsInspector) finish(a *Analyzer) error {
	return nil
}
