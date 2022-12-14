package nntp

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"

	"github.com/gorilla/websocket"
)

func isSilentError(err error) bool {
	return errors.Is(err, websocket.ErrCloseSent) ||
		errors.Is(err, net.ErrClosed) ||
		websocket.IsCloseError(err, websocket.CloseAbnormalClosure)
}

func (c *Client) readFromWSLoop() {
	var remoteReader bufio.Reader
	var remoteWriter bufio.Writer

	defer c.cancel()

	for {
		select {
		case <-c.ctx.Done():
			log.Println("finishing read ws loop")
			return
		default:
			msgType, r, err := c.wsConn.NextReader()
			if err != nil {
				if !isSilentError(err) {
					log.Println("error fetching from reader:", err)
				}
				return
			}

			if msgType != websocket.TextMessage {
				log.Println("received message of type", msgTypeString(msgType), "; ignoring")
				continue
			}

			remoteReader.Reset(r)
			line, err := remoteReader.ReadString('\n')
			for ; err == nil; line, err = remoteReader.ReadString('\n') {
				remoteWriter.Reset(c.remoteConn)
				_, writeErr := remoteWriter.WriteString(line)
				if writeErr != nil {
					log.Println("error writing line to remote connection:", writeErr)
					continue
				}
				flushErr := remoteWriter.Flush()
				if flushErr != nil {
					log.Println("error flushing remote writer:", flushErr)
				}
			}

			if !errors.Is(err, io.EOF) {
				log.Println("read loop encountered non EOF error:", err.Error())
			}

		}
	}
}
