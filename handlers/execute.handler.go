package handler

import (
	"fmt"
	"net/http"

	service "../services"
	types "../types"
	"github.com/google/uuid"
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

		//Setup directory
		service.CreateDirectory(dir)

		//Create files
		service.CopyExecuteJs(dir)
		service.CreateCodeJs(dir, _b.Code)
		service.CreateScriptSh(dir, service.StartSh)

		//roll up container and watch for file changes
		service.RollUpContiner(dir)

		w.WriteHeader(http.StatusOK)
		resp.status = true
		resp.message = "response"
		fmt.Fprintln(w, resp)
	})
}
