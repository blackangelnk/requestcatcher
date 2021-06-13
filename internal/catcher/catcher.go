package catcher

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/blackangelnk/requestcatcher/internal/config"
	"github.com/blackangelnk/requestcatcher/internal/request"
)

type Catcher struct {
	Server    *http.Server
	Broadcast chan *request.CaughtRequest
}

func NewCatcher(cfg *config.Configuration) *Catcher {
	mux := http.NewServeMux()
	c := &Catcher{
		Server: &http.Server{
			Addr:    ":" + strconv.Itoa(cfg.CatcherPort),
			Handler: mux,
		},
		Broadcast: make(chan *request.CaughtRequest),
	}
	mux.HandleFunc("/", c.handler)
	return c
}

func (c *Catcher) handler(w http.ResponseWriter, r *http.Request) {
	cr, err := c.Catch(r)
	if err != nil {
		log.Fatal("Failed to catch request", err)
	} else {
		c.Broadcast <- cr
	}
}

func (c *Catcher) Catch(r *http.Request) (*request.CaughtRequest, error) {
	headers, err := json.Marshal(r.Header)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	cr := &request.CaughtRequest{
		Url:           r.URL.String(),
		Time:          time.Now(),
		Method:        r.Method,
		ContentLength: r.ContentLength,
		RemoteAddr:    r.RemoteAddr,
		Body:          string(body),
		Headers:       string(headers),
	}
	return cr, nil
}

func (c *Catcher) Run() error {
	log.Printf("Running catcher")
	return c.Server.ListenAndServe()
}

func (c *Catcher) Stop(ctx context.Context) error {
	return c.Server.Shutdown(ctx)
}
