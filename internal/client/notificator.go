package client

import (
	"github.com/blackangelnk/requestcatcher/internal/request"
	"github.com/gorilla/websocket"
)

type notificator struct {
	wsClients map[*wsClient]struct{}
	register  chan *websocket.Conn
	delete    chan *wsClient
	Send      chan *request.CaughtRequest
}

func (n *notificator) Run() {
	for {
		select {
		case conn := <-n.register:
			c := &wsClient{
				conn: conn,
				send: make(chan *request.CaughtRequest),
				n:    n,
			}
			n.wsClients[c] = struct{}{}
			c.run()
		case c := <-n.delete:
			if _, ok := n.wsClients[c]; ok {
				delete(n.wsClients, c)
				close(c.send)
			}
		case cr := <-n.Send:
			for c := range n.wsClients {
				c.send <- cr
			}
		}
	}
}
