package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/prashant-raghu/computeEngine/types"
)

type Resp struct {
	status  bool
	message string
}

func Execute() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//create temp dir with code.js main.js script.sh//
		//roll up container
		//set timeout
		//add watcher of /out.txt
		var _b types.Execute
		var resp Resp
		r.ParseForm()
		_b.Code = r.FormValue("code")
		dir := uuid.New()
		_, err := os.Stat(fmt.Sprintf("temp/%s", dir.String()))
		if os.IsNotExist(err) {
			errDir := os.MkdirAll("temp/"+dir.String(), 0755)
			if errDir != nil {
				log.Fatal(err)
			}
		}
		b, err := ioutil.ReadFile(fmt.Sprintf("%s", "execute.js"))
		err = ioutil.WriteFile(fmt.Sprintf("temp/%s/%s", dir.String(), "execute.js"), b, 0644)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(fmt.Sprintf("temp/%s/%s", dir.String(), "code.js"), []byte(_b.Code), 0644)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		resp.status = true
		resp.message = "response"
		fmt.Fprintln(w, resp)
	})
}
