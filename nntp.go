package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type nntpHandler struct {
	upgrader     websocket.Upgrader
	connStateMap sync.Map
}

func NewNNTPHandler() nntpHandler {
	return nntpHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		connStateMap: sync.Map{},
	}
}

func (n nntpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := getAuthToken(r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("error authenticating request: %v", err), http.StatusUnauthorized,
		)
		return
	}
	_, err = n.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
