package client

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/blackangelnk/requestcatcher/internal/catcher"
	"github.com/blackangelnk/requestcatcher/internal/config"
	"github.com/jmoiron/sqlx"
)

type Client struct {
	Server *http.Server
	db     *sqlx.DB
}

type VRequest struct {
	catcher.CaughtRequest
}

func (v *VRequest) ParsedHeaders() map[string][]string {
	var headers map[string][]string
	err := json.Unmarshal([]byte(v.Headers), &headers)
	if err != nil {
		log.Print("Failed to unmarshal headers json", err)
	}
	return headers
}

func NewClient(cfg *config.Configuration, db *sqlx.DB) *Client {
	mux := http.NewServeMux()
	c := &Client{
		db: db,
		Server: &http.Server{
			Addr:    ":" + strconv.Itoa(cfg.ClientPort),
			Handler: mux,
		},
	}
	mux.HandleFunc("/", c.handler)
	return c
}

func (c *Client) handler(w http.ResponseWriter, r *http.Request) {
	var requests []VRequest
	err := c.db.Select(&requests, "SELECT * FROM request ORDER BY id DESC;")
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
