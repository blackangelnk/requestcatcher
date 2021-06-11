package catcher

import (
	"net/http"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/blackangelnk/requestcatcher/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type mReq struct {
	url         string
	method      string
	body        string
	len         int
	remote_addr string
	headers     map[string][]string
}

func TestCreate(t *testing.T) {
	db, _ := createDB(t)
	defer db.Close()
	cfg := &config.Configuration{
		ClientPort:  8077,
		CatcherPort: 8076,
	}
	c := NewCatcher(cfg, db)

	assert.NotNil(t, c)
}

func TestCatch(t *testing.T) {
	db, mock := createDB(t)
	defer db.Close()
	cfg := &config.Configuration{
		ClientPort:  8077,
		CatcherPort: 8076,
	}
	c := NewCatcher(cfg, db)
	cases := []mReq{
		{
			url:         "/case-1",
			method:      "POST",
			body:        "body",
			len:         4,
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
			len:         0,
			remote_addr: "127.0.0.1:38271",
			headers: map[string][]string{
				"User-Agent": {"curl"},
				"Accept":     {"*/*"},
			},
		},
	}
	for i, cs := range cases {
		mock.ExpectExec("INSERT INTO request ").WillReturnResult(sqlmock.NewResult(int64(i), 1))
		req, _ := http.NewRequest(cs.method, cs.url, strings.NewReader(cs.body))
		req.RemoteAddr = cs.remote_addr
		req.Header = cs.headers
		cr, err := c.Catch(req)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, int64(i), cr.Id)
	}
}

func createDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	return sqlx.NewDb(db, "sqlmock"), mock
}
