package fileIO

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	UA        string     `yaml:"ua"`
	Platforms [][]string `yaml:"platforms"`
}

func ReadConfig() (*Config, int) {
	file, err := os.ReadFile("./Source.yaml")
	if err != nil {
		fmt.Println(err)
	}

	cf := Config{}

	yaml.Unmarshal(file, &cf)
	L := len(cf.Platforms)
	return &cf, L
}
