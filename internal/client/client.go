package client

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/blackangelnk/requestcatcher/internal/config"
	"github.com/blackangelnk/requestcatcher/internal/request"
	"github.com/blackangelnk/requestcatcher/internal/storage"
	"github.com/gorilla/websocket"
)

type Client struct {
	Server      *http.Server
	Storage     storage.Storage
	Notificator *notificator
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewClient(cfg *config.Configuration, s storage.Storage) *Client {
	mux := http.NewServeMux()
	c := &Client{
		Storage: s,
		Server: &http.Server{
			Addr:    ":" + strconv.Itoa(cfg.ClientPort),
			Handler: mux,
		},
		Notificator: &notificator{
			wsClients: make(map[*wsClient]struct{}),
			register:  make(chan *websocket.Conn),
			delete:    make(chan *wsClient),
			Send:      make(chan *request.CaughtRequest),
		},
	}
	mux.HandleFunc("/", c.handler)
	mux.HandleFunc("/ws", c.handleWs)
	return c
}

func (c *Client) handler(w http.ResponseWriter, r *http.Request) {
	requests, err := c.Storage.Get()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("Failed to get requests from storage", err)
		return
	}
	t, err := template.ParseFiles("../../web/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("Failed to parse template files", err)
		return
	}
	err = t.Execute(w, requests)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("Failed to execute template", err)
	}
}

func (c *Client) handleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Failed to upgrade ws connection", err)
		return
	}
	c.Notificator.register <- conn
}

func (c *Client) Run() error {
	log.Printf("Running client")
	go c.Notificator.Run()
	return c.Server.ListenAndServe()
}

func (c *Client) Stop(ctx context.Context) error {
	return c.Server.Shutdown(ctx)
}
