package fileIO

import (
	"io/ioutil"
	"log"
	"os"
)

func Read_tmp() string {
	file, err := os.Open("./README_TEMP.md")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	return string(content)
}
