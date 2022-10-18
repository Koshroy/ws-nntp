package nntp

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

const PingTimeout = 60 * 5
const ReadLoopBufSize = 4096
const BadConnectLine = "550 could not connect to remote"

type Client struct {
	wsConn                *websocket.Conn
	remoteConn            net.Conn
	remoteConnEstablished bool
	ctx                   context.Context
	pingTimer             *time.Ticker
	cancel                context.CancelFunc
}

func (c *Client) pongHandler(_ string) error {
	c.pingTimer.Reset(PingTimeout * time.Second)
	return nil
}

func newClient(wsConn *websocket.Conn, remoteConn net.Conn, ctx context.Context) *Client {
	ticker := time.NewTicker(PingTimeout * time.Second)
	newCtx, cancel := context.WithCancel(ctx)
	client := &Client{
		wsConn:                wsConn,
		remoteConn:            remoteConn,
		remoteConnEstablished: false,
		ctx:                   newCtx,
		pingTimer:             ticker,
		cancel:                cancel,
	}
	client.wsConn.SetPongHandler(client.pongHandler)
	return client
}

func (c *Client) managerLoop() {
	ctxDone := c.ctx.Done()
	closeMsg := websocket.FormatCloseMessage(1000, "closing websocket")
	for {
		select {
		case <-ctxDone:
			err := c.wsConn.WriteMessage(websocket.CloseMessage, closeMsg)
			if err != nil && err != websocket.ErrCloseSent {
				log.Println("error writing close message to websocket connection:", err)
			}
			err = c.wsConn.Close()
			if err != nil {
				log.Println("error closing websocket connection:", err)
			}
			err = c.remoteConn.Close()
			if err != nil {
				log.Println("error closing remote nntp connection:", err)
			}
			return
		case <-c.pingTimer.C:
			log.Println("ping timeout reached, closing connection")
			err := c.wsConn.Close()
			if err != nil {
				log.Println("error closing websocket connection:", err)
			}
			err = c.remoteConn.Close()
			if err != nil {
				log.Println("error closing remote nntp connection:", err)
			}
			return
		}
	}
}

func msgTypeString(msgType int) string {
	switch msgType {
	case websocket.TextMessage:
		return "text"
	case websocket.BinaryMessage:
		return "binary"
	case websocket.PingMessage:
		return "ping"
	case websocket.PongMessage:
		return "pong"
	case websocket.CloseMessage:
		return "close"
	default:
		return "unknown"
	}
}
