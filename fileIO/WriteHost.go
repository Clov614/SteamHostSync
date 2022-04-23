package fileIO

import (
	"io/ioutil"
)

func WriteHost(content string, filename string) {
	ioutil.WriteFile("./"+filename, []byte(content), 0666)
}
