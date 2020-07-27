package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	// "types"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type key int

type Execute struct {
	// Age  int
	code string
}

type Resp struct {
	status  bool
	message string
}

const (
	requestIDKey key = 0
)

var (
	Version      string = ""
	GitTag       string = ""
	GitCommit    string = ""
	GitTreeState string = ""
	listenAddr   string
	healthy      int32
)

func main() {
	flag.StringVar(&listenAddr, "listen-addr", ":5000", "server listen address")
	flag.Parse()

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	logger.Println("Simple go server")
	logger.Println("Version:", Version)
	logger.Println("GitTag:", GitTag)
	logger.Println("GitCommit:", GitCommit)
	logger.Println("GitTreeState:", GitTreeState)

	logger.Println("Server is starting...")

	router := http.NewServeMux()
	router.Handle("/", index())
	router.Handle("/test", TodoShow())
	router.Handle("/healthz", healthz())
	router.Handle("/execute", execute())

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      tracing(nextRequestID)(logging(logger)(router)),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Println("Server is ready to handle requests at", listenAddr)
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Println("Server stopped")
}

func TodoShow() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		todoId := vars["todoId"]
		fmt.Fprintln(w, "Todo show:", todoId)
	})
}
func index() http.Handler {

	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	r.ParseForm()
	// 	code := r.FormValue("code")
	// 	fmt.Fprintf(w, code)
	// })
	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// 	fmt.Println(r.Body)
	// 	decoder := json.NewDecoder(r.Body)
	// 	var data Execute
	// 	err := decoder.Decode(&data)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	code := data.code
	// 	fmt.Println(r)
	// 	fmt.Println(data)
	// 	fmt.Println(code)
	// 	fmt.Fprintf(w, code)
	// })

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Server up and running")
		// fmt.Fprintln(w, _b.code)
	})
}

func execute() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//create temp dir with code.js main.js script.sh//
		//roll up container
		//set timeout
		//add watcher of /out.txt
		var _b Execute
		var resp Resp
		r.ParseForm()
		_b.code = r.FormValue("code")
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
		err = ioutil.WriteFile(fmt.Sprintf("temp/%s/%s", dir.String(), "code.js"), []byte(_b.code), 0644)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		resp.status = true
		resp.message = "response"
		fmt.Fprintln(w, resp)
	})
}

func healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&healthy) == 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}
func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
