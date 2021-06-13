package storage

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/blackangelnk/requestcatcher/internal/request"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	db, mock := createDB(t)
	defer db.Close()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS request").WillReturnResult(sqlmock.NewResult(0, 0))
	s := NewDb(db)

	assert.NotNil(t, s)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSave(t *testing.T) {
	db, mock := createDB(t)
	defer db.Close()

	mock.ExpectExec("INSERT INTO request ").WillReturnResult(sqlmock.NewResult(1, 1))
	s := &dbStorage{
		db: db,
	}
	cr := request.CaughtRequest{}
	_, err := s.Save(cr)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGet(t *testing.T) {
	db, mock := createDB(t)
	defer db.Close()

	cr := request.CaughtRequest{
		Id:            1,
		Time:          time.Now(),
		Method:        "GET",
		ContentLength: 0,
		RemoteAddr:    "127.0.0.1:123",
		Url:           "/",
		Body:          "",
	}
	rows := sqlmock.NewRows(
		[]string{"id", "created_at", "url", "headers", "body", "method", "remote_addr", "content_length"},
	).AddRow(cr.Id, cr.Time, cr.Url, "", cr.Body, cr.Method, cr.RemoteAddr, cr.ContentLength)

	mock.ExpectQuery("SELECT (.+) FROM request(.+)").WillReturnRows(rows)
	s := &dbStorage{
		db: db,
	}
	requests, err := s.Get()
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Len(t, requests, 1)
	assert.Equal(t, cr.Id, requests[0].Id)
}

func createDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	return sqlx.NewDb(db, "sqlmock"), mock
}
