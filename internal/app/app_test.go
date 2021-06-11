package app

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/blackangelnk/requestcatcher/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCatcher(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS request").WillReturnResult(sqlmock.NewResult(0, 0))
	cfg := &config.Configuration{}
	app := NewApp(cfg, sqlx.NewDb(db, "sqlmock"))
	assert.NotNil(t, app)
	assert.NoError(t, mock.ExpectationsWereMet())
}
