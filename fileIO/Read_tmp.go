package fileIO

import (
	"io/ioutil"
	"os"
)

func Read_tmp() string {
	pwd, _ := os.Getwd()
	file, err := os.Open(pwd + "\\README_TEMP.md")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	return string(content)
}
