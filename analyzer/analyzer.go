package analyzer

import (
	"crypto/md5"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"nornir/annotations"
	"nornir/io"
)

type Analyzer struct {
	output          string
	prefix          string
	restlib         string
	dir             *io.Directory
	ModName         string
	Structs         map[string]Struct
	StructImportMap map[string][]Import
	AnnotatedFiles  []*annotations.AnnotatedFile
	AnnotatedFuncs  []*annotations.AnnotatedFunc
	AnnotatedVars   []*annotations.AnnotatedVar
}

func NewAnalyzer(modname string, dir *io.Directory, outdir string, prefix string, restlib string) *Analyzer {
	return &Analyzer{
		output:          outdir,
		prefix:          prefix,
		restlib:         restlib,
		ModName:         modname,
		dir:             dir,
		Structs:         make(map[string]Struct),
		StructImportMap: make(map[string][]Import),
		AnnotatedFiles:  make([]*annotations.AnnotatedFile, 0),
		AnnotatedFuncs:  make([]*annotations.AnnotatedFunc, 0),
		AnnotatedVars:   make([]*annotations.AnnotatedVar, 0),
	}
}

func (a *Analyzer) AnalyzeStructs() error {
	fset := token.NewFileSet()
	inspector := &StructInspector{}

	err := a.parseDirectoriesRecursively(a.ModName, a.dir, fset, inspector)
	if err != nil {
		return err
	}

	return inspector.finish(a)
}

func (a *Analyzer) AnalyzeAnnotations() error {
	fset := token.NewFileSet()
	inspector := &AnnotationsInspector{}

	err := a.parseDirectoriesRecursively(a.ModName, a.dir, fset, inspector)
	if err != nil {
		return err
	}

	return inspector.finish(a)
}

func (a *Analyzer) parseFile(currentPath string, filename string, fset *token.FileSet, insepctor inspector) error {
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	log.Printf("Parsing file: %s", filename)

	ctx := &InspectorContext{
		Analyzer: a,
		Path:     currentPath,
		FileId:   generateFileId(filename),
		Filename: filename,
		Data:     make(map[string]interface{}),
	}

	ast.Inspect(file, func(node ast.Node) bool {
		return insepctor.visit(ctx, node)
	})

	return nil
}

func (a *Analyzer) parseDirectoriesRecursively(currentPath string, dir *io.Directory, fset *token.FileSet, insepctor inspector) error {
	for _, file := range dir.Files {
		err := a.parseFile(currentPath, file, fset, insepctor)
		if err != nil {
			return err
		}
	}

	for _, innerDir := range dir.Children {
		err := a.parseDirectoriesRecursively(currentPath+"/"+innerDir.Name, innerDir, fset, insepctor)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateFileId(filename string) string {
	md5sum := md5.Sum([]byte(filename))
	return string(md5sum[:])
}
