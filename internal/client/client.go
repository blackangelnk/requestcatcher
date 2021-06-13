package client

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/blackangelnk/requestcatcher/internal/config"
	"github.com/blackangelnk/requestcatcher/internal/storage"
)

type Client struct {
	Server  *http.Server
	Storage storage.Storage
}

func NewClient(cfg *config.Configuration, s storage.Storage) *Client {
	mux := http.NewServeMux()
	c := &Client{
		Storage: s,
		Server: &http.Server{
			Addr:    ":" + strconv.Itoa(cfg.ClientPort),
			Handler: mux,
		},
	}
	mux.HandleFunc("/", c.handler)
	return c
}

func (c *Client) handler(w http.ResponseWriter, r *http.Request) {
	requests, err := c.Storage.Get()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("Failed to select requests from database", err)
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

func (c *Client) Run() error {
	log.Printf("Running client")
	return c.Server.ListenAndServe()
}

func (c *Client) Stop(ctx context.Context) error {
	return c.Server.Shutdown(ctx)
}
