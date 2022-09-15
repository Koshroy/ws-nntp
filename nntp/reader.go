package nntp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

func (c *Client) readLoop() {
	var readBuf bytes.Buffer
	var remoteReader bufio.Reader
	var remoteWriter bufio.Writer

	for {
		msgType, r, err := c.wsConn.NextReader()
		if err != nil {
			if !errors.Is(err, websocket.ErrCloseSent) {
				log.Println("error fetching from reader:", err)
			}
			return
		}

		if msgType != websocket.TextMessage {
			log.Println("received message of type", msgTypeString(msgType), "; ignoring")
			continue
		}

		n, err := readBuf.ReadFrom(r)
		if err != nil {
			log.Printf("error reading websocket message of length %d: %v", n, err)
			continue
		}

		line, err := readBuf.ReadString('\n')
		if err != nil {
			log.Println("error reading string from read buffer:", err)
		}

		line = strings.Trim(line, "\n")

		if isSpecialCmd(line) {
			proc := strings.SplitN(line, " ", 2)
			if len(proc) < 2 {
				c.writeChan <- fmt.Sprintf("%s: malformed connect line", BadConnectLine)
				continue
			}

			remoteURL, err := url.Parse(strings.ToLower(proc[1]))
			if err != nil {
				c.writeChan <- fmt.Sprintf("%s: %v", BadConnectLine, err)
				continue
			}

			if remoteURL.Scheme != "nntp" && remoteURL.Scheme != "nntps" {
				c.writeChan <- fmt.Sprintf("%s: %v", BadConnectLine, err)
				continue
			}

			conn, err := net.Dial("tcp", remoteURL.Host)
			if err != nil {
				c.writeChan <- fmt.Sprintf("%s: %v", BadConnectLine, err)
				continue
			}

			remoteReader.Reset(conn)
			line, err = remoteReader.ReadString('\n')
			if err != nil {
				c.writeChan <- fmt.Sprintf("%s: %v", BadConnectLine, err)
				conn.Close()
				continue
			}

			c.writeChan <- line
			c.remoteConn = conn
			c.remoteConnEstablished = true
		} else {
			if !c.remoteConnEstablished {
				log.Println("tried to write bytes to unestablished remote connection")
				continue
			}

			remoteWriter.Reset(c.remoteConn)
			_, err := remoteWriter.WriteString(line)
			if err != nil {
				log.Println("error writing line to remote connection:", err)
				continue
			}

			_, err = readBuf.ReadFrom(c.remoteConn)
			if err != nil {
				log.Println("error reading from remote connection to buf:", err)
				continue
			}

			c.writeChan <- readBuf.String()
		}
	}
}
