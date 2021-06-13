package catcher

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blackangelnk/requestcatcher/internal/config"
	"github.com/stretchr/testify/assert"
)

type mReq struct {
	url         string
	method      string
	body        string
	remote_addr string
	headers     map[string][]string
}

var cases []mReq

func init() {
	cases = []mReq{
		{
			url:         "/case-1",
			method:      "POST",
			body:        "body",
			remote_addr: "127.0.0.1:38270",
			headers: map[string][]string{
				"Content-Type": {"application/json"},
				"User-Agent":   {"curl"},
			},
		},
		{
			url:         "/case-2",
			method:      "GET",
			body:        "",
			remote_addr: "127.0.0.1:38271",
			headers: map[string][]string{
				"User-Agent": {"curl"},
				"Accept":     {"*/*"},
			},
		},
	}
}

func TestCreate(t *testing.T) {
	cfg := &config.Configuration{
		ClientPort:  8077,
		CatcherPort: 8076,
	}
	c := NewCatcher(cfg)

	assert.NotNil(t, c)
}

func TestCatch(t *testing.T) {
	cfg := &config.Configuration{
		ClientPort:  8077,
		CatcherPort: 8076,
	}
	c := NewCatcher(cfg)

	for _, cs := range cases {
		req, _ := http.NewRequest(cs.method, cs.url, strings.NewReader(cs.body))
		req.RemoteAddr = cs.remote_addr
		req.Header = cs.headers
		cr, err := c.Catch(req)
		assert.NoError(t, err)
		assert.Equal(t, cr.Url, cs.url)
		assert.Equal(t, cr.Method, cs.method)
		assert.Equal(t, cr.Body, cs.body)
		assert.Equal(t, cr.ContentLength, int64(len(cs.body)))
		assert.Equal(t, cr.RemoteAddr, cs.remote_addr)
		headers, _ := json.Marshal(cs.headers)
		assert.Equal(t, cr.Headers, string(headers))
	}
}

func TestHandler(t *testing.T) {
	cfg := &config.Configuration{
		ClientPort:  8077,
		CatcherPort: 8076,
	}
	c := NewCatcher(cfg)

	for _, cs := range cases {
		req, _ := http.NewRequest(cs.method, cs.url, strings.NewReader(cs.body))
		req.RemoteAddr = cs.remote_addr
		req.Header = cs.headers
		w := httptest.NewRecorder()
		go c.handler(w, req)
		cr := <-c.Broadcast
		assert.Equal(t, cr.Url, cs.url)
		assert.Equal(t, cr.Method, cs.method)
		assert.Equal(t, cr.Body, cs.body)
		assert.Equal(t, cr.ContentLength, int64(len(cs.body)))
		assert.Equal(t, cr.RemoteAddr, cs.remote_addr)
		headers, _ := json.Marshal(cs.headers)
		assert.Equal(t, cr.Headers, string(headers))
	}
}
