package io

import (
	"bufio"
	"os"
	"strings"
)

func ReadModuleName() string {
	gomod, err := os.Open("go.mod")
	if err != nil {
		panic(err)
	}

	defer gomod.Close()

	info, err := gomod.Stat()
	if err != nil {
		panic(err)
	}

	if info.Size() == 0 {
		panic("go.mod is empty")
	}

	scanner := bufio.NewScanner(gomod)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "module") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}

	panic("module name not found")
}
