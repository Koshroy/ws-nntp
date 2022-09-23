package nntp

import (
	"bufio"
	"errors"
	"io"
	"log"

	"github.com/gorilla/websocket"
)

func (c *Client) writeToWSLoop() {
	var remoteReader bufio.Reader
	var wsWriter bufio.Writer

	for {
		remoteReader.Reset(c.remoteConn)
		line, err := remoteReader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("Received line", line, "EOF encountered in remote stream. Exiting")
				return
			} else {
				log.Println("Received error:", err)
			}
		}

		w, err := c.wsConn.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Println("error getting writer from websocket connection:", err)
			continue
		}

		wsWriter.Reset(w)
		wsWriter.WriteString(line)
		wsWriter.Flush()
		err = w.Close()
		if err != nil {
			log.Println("error finishing write to websocket connection:", err)
		}
	}
}
