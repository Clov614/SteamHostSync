package fileIO

import (
	"io/ioutil"
	"os"
)

func WriteHost(content string, filename string) {
	pwd, _ := os.Getwd()
	ioutil.WriteFile(pwd+"\\"+filename, []byte(content), 0666)
}
