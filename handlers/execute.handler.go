package handler

import (
	"fmt"
	"log"
	"net/http"

	service "../services"
	types "../types"

	"github.com/google/uuid"
	"github.com/yaacov/observer/observer"
)

type Resp struct {
	status  bool
	message string
}

func Execute() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var _b types.Execute
		var resp Resp
		r.ParseForm()
		_b.Code = r.FormValue("code")
		dir := uuid.New()
		var toWatch string = fmt.Sprintf("./%s/%s/out.txt", service.ParentDir, dir.String())

		//Setup directory
		service.CreateDirectory(dir)

		//Create files
		service.CopyExecuteJs(dir)
		service.CreateCodeJs(dir, _b.Code)
		service.CreateScriptSh(dir, service.StartSh)

		//Add Watcher for file create
		o := observer.Observer{}
		err := o.Watch([]string{toWatch})
		if err != nil {
			log.Fatal("Error: ", err)
		}
		defer o.Close()
		o.AddListener(func(e interface{}) {
			if e.(observer.WatchEvent).Op == 2 {
				w.WriteHeader(http.StatusOK)
				resp.status = true
				resp.message = service.RetrieveOutTxt(dir)
				fmt.Fprintln(w, resp.message)
			}
		})

		//roll up container
		service.RollUpContiner(dir)
	})
}
