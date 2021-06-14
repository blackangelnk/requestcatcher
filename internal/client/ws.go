package client

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/blackangelnk/requestcatcher/internal/request"
	"github.com/gorilla/websocket"
)

const (
	pingPeriod = 30 * time.Second
	writeWait  = 5 * time.Second
	pongWait   = 30 * time.Second
	readLimit  = 256
)

type wsClient struct {
	conn   *websocket.Conn
	send   chan *request.CaughtRequest
	n      *notificator
	closer sync.Once
}

func (c *wsClient) run() {
	log.Print("Running new ws client")
	go c.read()
	go c.write()
}

func (c *wsClient) close() {
	c.closer.Do(func() {
		log.Print("Closing ws client")
		c.n.delete <- c
		c.conn.SetWriteDeadline(time.Now().Add(writeWait))
		c.conn.WriteMessage(websocket.CloseMessage, []byte{})

		c.conn.Close()
	})
}

func (c *wsClient) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.close()
	}()
	for {
		select {
		case cr, ok := <-c.send:
			if !ok {
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			msg, _ := json.Marshal(cr)
			w.Write(msg)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *wsClient) read() {
	defer c.close()
	c.conn.SetReadLimit(readLimit)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(msg string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		if _, _, err := c.conn.NextReader(); err != nil {
			break
		}
	}
}
