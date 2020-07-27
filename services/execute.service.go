package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/uuid"
)

const (
	parentDir string = "temp"
)

func CreateDirectory(dir uuid.UUID) {
	_, err := os.Stat(fmt.Sprintf("%s/%s", parentDir, dir.String()))
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(fmt.Sprintf("%s/%s", parentDir, dir.String()), 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}
}

func CopyExecuteJs(dir uuid.UUID) {
	b, err := ioutil.ReadFile(fmt.Sprintf("%s", "./assets/execute.js"))
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s/%s", parentDir, dir.String(), "execute.js"), b, 0644)
	if err != nil {
		panic(err)
	}
}

func CreateCodeJs(dir uuid.UUID, code string) {
	err := ioutil.WriteFile(fmt.Sprintf("%s/%s/%s", parentDir, dir.String(), "code.js"), []byte(code), 0644)
	if err != nil {
		panic(err)
	}
}

func CreateScriptSh(dir uuid.UUID, content string) {
	err := ioutil.WriteFile(fmt.Sprintf("%s/%s/%s", parentDir, dir.String(), "script.sh"), []byte(content), 0644)
	if err != nil {
		panic(err)
	}
}
