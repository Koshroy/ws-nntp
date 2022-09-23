package nntp

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"

	"github.com/gorilla/websocket"
)

func (c *Client) writeToWSLoop() {
	var remoteReader bufio.Reader
	var wsWriter bufio.Writer

	defer c.cancel()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			remoteReader.Reset(c.remoteConn)
			line, err := remoteReader.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Println("Received line", line, "EOF encountered in remote stream. Exiting")
					return
				} else if errors.Is(err, net.ErrClosed) {
					return
				} else {
					log.Println("Received line", line, "with error:", err)
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
}
