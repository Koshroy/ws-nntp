package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	nntp := nntpHandler{}
	r.Path("/nntp").Handler(nntp)
	http.Handle("/", r)
}
