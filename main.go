package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Koshroy/ws-nntp/nntp"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	nntp := nntp.Handler{}
	r.Path("/nntp").Handler(nntp)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:9090",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("starting web server")
	log.Fatal(srv.ListenAndServe())
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
