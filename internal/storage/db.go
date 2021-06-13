package storage

import (
	"log"

	"github.com/blackangelnk/requestcatcher/internal/request"
	"github.com/jmoiron/sqlx"
)

type dbStorage struct {
	db *sqlx.DB
}

func NewDb(db *sqlx.DB) *dbStorage {
	s := &dbStorage{
		db: db,
	}
	s.init()
	return s
}

func (s *dbStorage) init() {
	s.db.MustExec(`CREATE TABLE IF NOT EXISTS request (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"url" TEXT,
		"headers" TEXT,
		"body" TEXT,
		"method" TEXT,
		"remote_addr" TEXT,
		"content_length" INTEGER
	)`)
}

func (s *dbStorage) Save(req request.CaughtRequest) (request.CaughtRequest, error) {
	q := `INSERT INTO request (url, method, content_length, remote_addr, headers, body)
	 VALUES(:url,:method,:content_length,:remote_addr,:headers,:body)`

	res, err := s.db.NamedExec(q, req)
	if err != nil {
		return req, err
	}
	req.Id, _ = res.LastInsertId()
	return req, nil
}

func (s *dbStorage) Get() ([]request.CaughtRequest, error) {
	var requests []request.CaughtRequest
	err := s.db.Select(&requests, "SELECT * FROM request ORDER BY id DESC;")
	if err != nil {
		log.Print("Failed to select requests from database", err)
	}
	return requests, err
}
