package nntp

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Handler struct {
	upgrader     websocket.Upgrader
	connStateMap sync.Map
}

func NewNNTPHandler() Handler {
	return Handler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		connStateMap: sync.Map{},
	}
}

func (n Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := n.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
