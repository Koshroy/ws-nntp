package nntp

import (
	"context"
	"log"
	"net"
	"strings"
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
	writeChan             chan string
}

func (c *Client) pongHandler(_ string) error {
	c.pingTimer.Reset(PingTimeout * time.Second)
	return nil
}

func newClient(wsConn *websocket.Conn, ctx context.Context) *Client {
	writeChan := make(chan string)
	ticker := time.NewTicker(PingTimeout * time.Second)
	client := &Client{
		wsConn:                wsConn,
		remoteConn:            nil,
		remoteConnEstablished: false,
		ctx:                   ctx,
		pingTimer:             ticker,
		writeChan:             writeChan,
	}
	client.wsConn.SetPongHandler(client.pongHandler)
	return client
}

func (c *Client) managerLoop() {
	ctxDone := c.ctx.Done()
	for {
		select {
		case <-ctxDone:
			log.Println("Finishing websocket read/write loops")
			return
		case <-c.pingTimer.C:
			log.Println("ping timeout reached, closing connection")
			close(c.writeChan)
			err := c.wsConn.Close()
			if err != nil {
				log.Println("error closing websocket connection:", err)
			}
			if c.remoteConnEstablished {
				err := c.remoteConn.Close()
				if err != nil {
					log.Println("error closing remote nntp connection:", err)
				}
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

func isSpecialCmd(line string) bool {
	lowerLine := strings.ToLower(line)
	return strings.HasPrefix(lowerLine, "connect ")
}
