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
