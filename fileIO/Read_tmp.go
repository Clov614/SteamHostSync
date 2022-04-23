package fileIO

import (
	"io/ioutil"
	"os"
)

func Read_tmp() string {
	file, err := os.Open("./README_TEMP.md")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	return string(content)
}
