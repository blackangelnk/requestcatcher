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
	"github.com/jmoiron/sqlx"
)

type Catcher struct {
	db     *sqlx.DB
	Server *http.Server
}

type CaughtRequest struct {
	Id            int64
	Time          time.Time `db:"created_at"`
	Method        string    `db:"method"`
	ContentLength int64     `db:"content_length"`
	RemoteAddr    string    `db:"remote_addr"`
	Url           string    `db:"url"`
	Headers       string    `db:"headers"`
	Body          string    `db:"body"`
}

func NewCatcher(cfg *config.Configuration, db *sqlx.DB) *Catcher {
	mux := http.NewServeMux()
	c := &Catcher{
		db: db,
		Server: &http.Server{
			Addr:    ":" + strconv.Itoa(cfg.CatcherPort),
			Handler: mux,
		},
	}
	mux.HandleFunc("/", c.handler)
	return c
}

func (c *Catcher) handler(w http.ResponseWriter, r *http.Request) {
	_, err := c.Catch(r)
	if err != nil {
		log.Fatal("Failed to catch request", err)
	}
}

func (c *Catcher) Catch(r *http.Request) (*CaughtRequest, error) {
	q := `INSERT INTO request (url, method, content_length, remote_addr, headers, body)
	 VALUES(:url,:method,:content_length,:remote_addr,:headers,:body)`
	headers, err := json.Marshal(r.Header)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	cr := &CaughtRequest{
		Url:           r.URL.String(),
		Time:          time.Now(),
		Method:        r.Method,
		ContentLength: r.ContentLength,
		RemoteAddr:    r.RemoteAddr,
		Body:          string(body),
		Headers:       string(headers),
	}
	res, err := c.db.NamedExec(q, cr)
	if err != nil {
		return nil, err
	}
	cr.Id, _ = res.LastInsertId()
	return cr, nil
}

func (c *Catcher) Run() error {
	log.Printf("Running catcher")
	return c.Server.ListenAndServe()
}

func (c *Catcher) Stop(ctx context.Context) error {
	return c.Server.Shutdown(ctx)
}
