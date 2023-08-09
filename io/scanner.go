package io

import (
	"os"
	"path/filepath"
	"strings"
)

type Directory struct {
	Name     string
	Children []*Directory
	Files    []string
}

func FindFiles() *Directory {
	baseDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return iterateDirectory(&Directory{
		Name: filepath.Base(baseDir),
	}, baseDir)
}

func iterateDirectory(currentDir *Directory, path string) *Directory {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			innerDir := &Directory{
				Name: file.Name(),
			}

			currentDir.Children = append(currentDir.Children, innerDir)

			iterateDirectory(innerDir, path+"/"+file.Name())
		} else {
			if !strings.HasSuffix(file.Name(), ".go") {
				continue
			}

			currentDir.Files = append(currentDir.Files, path+"/"+file.Name())
		}
	}

	return currentDir
}
