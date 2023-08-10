package generator

import (
	"nornir/log"
	"os"
	"path"
)

func GenerateTypeScriptClient(usages []usage) {
	createFileOnce("client.ts", []string{
		"import axios from 'axios';",
		"",
		"const DOMAIN = 'http://localhost:8080';",
		"",
		"export const restClient = axios.create({",
		"  baseURL: DOMAIN,",
		"});",
	})

	/*imports := []string{"import { restClient } from './client';"}
	body := make([]string, 0)*/

}

func createFileOnce(filename string, content []string) {
	outfile := path.Join(cfg.OutputDir, filename)

	if _, err := os.Stat(outfile); os.IsNotExist(err) {
		log.Debugf("Generating %s", outfile)

		file, err := os.Create(outfile)
		if err != nil {
			panic(err)
		}

		defer file.Close()

		for _, line := range content {
			file.WriteString(line + "\n")
		}
	}
}
