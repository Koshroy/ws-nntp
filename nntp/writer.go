package nntp

import (
	"bytes"
	"log"

	"github.com/gorilla/websocket"
)

func (c *Client) writeLoop() {
	var writeBuf bytes.Buffer
	for line := range c.writeChan {
		writeBuf.WriteString(line)
		err := c.wsConn.WriteMessage(websocket.TextMessage, writeBuf.Bytes())
		if err != nil {
			log.Println("error writing to websocket:", err)
		}
	}
}
