package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

func CreateDirectory(dir uuid.UUID) {
	_, err := os.Stat(fmt.Sprintf("%s/%s", ParentDir, dir.String()))
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(fmt.Sprintf("%s/%s", ParentDir, dir.String()), 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}
}

func CopyExecuteJs(dir uuid.UUID) {
	b, err := ioutil.ReadFile(fmt.Sprintf("%s", "./assets/execute.js"))
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s/%s", ParentDir, dir.String(), "execute.js"), b, 0644)
	if err != nil {
		panic(err)
	}
}

func CreateCodeJs(dir uuid.UUID, code string) {
	err := ioutil.WriteFile(fmt.Sprintf("%s/%s/%s", ParentDir, dir.String(), "code.js"), []byte(code), 0644)
	if err != nil {
		panic(err)
	}
}

func CreateScriptSh(dir uuid.UUID, content string) {
	var scriptDir = fmt.Sprintf("%s/%s/%s", ParentDir, dir.String(), "script.sh")
	var scriptDockerDir = fmt.Sprintf("%s/%s/%s", ParentDir, dir.String(), "scriptDocker.sh")
	err := ioutil.WriteFile(scriptDir, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
	//CreateExecutable
	err = exec.Command("chmod", "+x", scriptDir).Run()
	if err != nil {
		panic(err)
	}
	var rollUpBash = fmt.Sprintf(" sudo docker run --name %s --mount type=bind,source=\"$(pwd)\"/%s/%s,target=/app sandbox:v1", dir.String(), ParentDir, dir.String())
	err = ioutil.WriteFile(scriptDockerDir, []byte(rollUpBash), 0644)
	err = exec.Command("chmod", "+x", scriptDockerDir).Run()
	if err != nil {
		panic(err)
	}
}

func RollUpContiner(dir uuid.UUID) {
	// sudo docker run --name a8b558cf-cacf-4898-9743-b0b02007c059 --mount type=bind,source="$(pwd)"/temp/a8b558cf-cacf-4898-9743-b0b02007c059,target=/app sandbox:v1
	var toExec string = fmt.Sprintf("./%s/%s/scriptDocker.sh", ParentDir, dir.String())

	cmd := exec.Command("/bin/sh", toExec)
	cmd.Run()
}

func RetrieveOutTxt(dir uuid.UUID) string {
	// sudo docker run --name a8b558cf-cacf-4898-9743-b0b02007c059 --mount type=bind,source="$(pwd)"/temp/a8b558cf-cacf-4898-9743-b0b02007c059,target=/app sandbox:v1
	var toRet string = fmt.Sprintf("./%s/%s/out.txt", ParentDir, dir.String())
	b, err := ioutil.ReadFile(toRet)
	if err != nil {
		panic(err)
	}
	return string(b)
}
