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
	pwd, _ := os.Getwd()
	file, err := os.ReadFile(pwd + "\\Source.yaml")
	if err != nil {
		fmt.Println(err)
	}

	cf := Config{}

	yaml.Unmarshal(file, &cf)
	L := len(cf.Platforms)
	return &cf, L
}
