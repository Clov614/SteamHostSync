package fileIO

import (
	"io/ioutil"
	"os"
)

func ReadHtml() string {
	file, err := os.Open("./hosts.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	return string(content)
}
