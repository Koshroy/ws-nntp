package main

import (
	"net/http"

	"github.com/Koshroy/ws-nntp/nntp"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	nntp := nntp.Handler{}
	r.Path("/nntp").Handler(nntp)
	http.Handle("/", r)
}
