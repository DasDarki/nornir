package generator

import (
	"nornir/io"
	"os"
)

var cfg *io.Config

func Prepare(conf *io.Config) {
	cfg = conf

	os.MkdirAll(cfg.OutputDir, os.ModePerm)
}
