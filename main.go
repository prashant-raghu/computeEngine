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

	"github.com/prashant-raghu/computeEngine/services"
	"./types"

	"github.com/google/uuid"
)

type key int

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
	router.Handle("/healthz", service.Healthz(healthy))
	router.Handle("/execute", execute())

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      (tracing(nextRequestID))(logging(logger)(router)),
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

func index() http.Handler {
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
	})
}

func execute() http.Handler {
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
