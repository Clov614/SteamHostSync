package fileIO

import "io/ioutil"

func WriteFile(content, fileName string) {
	err := ioutil.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
}
