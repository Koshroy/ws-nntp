package nntp

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
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
	host := r.Header.Get("remote")
	if host == "" {
		http.Error(w, "no remote header found", http.StatusBadRequest)
		return
	}

	nntpURL, err := url.Parse(host)
	if err != nil {
		errStr := fmt.Sprintf("error parsing remote as URL: %s", err.Error())
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}

	scheme := nntpURL.Scheme
	if scheme != "nntp" {
		http.Error(w, "only nntp scheme URLs are allowed", http.StatusBadRequest)
		return
	}

	nntpConn, err := net.Dial("tcp", nntpURL.Host)
	if err != nil {
		errStr := fmt.Sprintf("error opening connection to %s: %s", host, err.Error())
		http.Error(w, errStr, http.StatusInternalServerError)
		return
	}

	wsConn, err := n.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		closeErr := nntpConn.Close() != nil
		if closeErr {
			log.Println("error closing nntp connection:", closeErr)
		}
		return
	}

	client := newClient(wsConn, nntpConn, context.Background())

	go client.managerLoop()
	go client.writeToWSLoop()
	go client.readFromWSLoop()
}
