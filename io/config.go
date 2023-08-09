package io

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	RestLib   string `yaml:"rest_lib"`
	OutputDir string `yaml:"output_dir"`
	Prefix    string `yaml:"prefix"`
}

func (cfg *Config) Save() {
	file, err := os.Create("nornir.yml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	if err := encoder.Encode(cfg); err != nil {
		panic(err)
	}
}

func LoadConfig() *Config {
	file, err := os.Open("nornir.yml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	cfg := &Config{}
	if err := decoder.Decode(cfg); err != nil {
		panic(err)
	}

	return cfg
}

func IsConfigExists() bool {
	if _, err := os.Stat("negroni.yml"); os.IsNotExist(err) {
		return false
	}

	return true
}
