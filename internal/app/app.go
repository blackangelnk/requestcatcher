package app

import (
	"context"
	"log"
	"net/http"

	"github.com/blackangelnk/requestcatcher/internal/catcher"
	"github.com/blackangelnk/requestcatcher/internal/client"
	"github.com/blackangelnk/requestcatcher/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	Client  *client.Client
	DB      *sqlx.DB
	Catcher *catcher.Catcher
}

func NewApp(cfg *config.Configuration, db *sqlx.DB) *App {
	app := &App{
		Client:  client.NewClient(cfg, db),
		DB:      db,
		Catcher: catcher.NewCatcher(cfg, db),
	}
	app.initDB()
	return app
}

func (a *App) initDB() {
	a.DB.MustExec(`CREATE TABLE IF NOT EXISTS request (
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

func (a *App) Run() {
	go func() {
		err := a.Catcher.Run()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Error while running catcher", err)
			panic(err)
		}
	}()
	go func() {
		err := a.Client.Run()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Error while running client", err)
			panic(err)
		}
	}()
}

func (a *App) Stop(ctx context.Context) {
	err := a.Catcher.Stop(ctx)
	if err != nil {
		log.Fatal("Error while stopping catcher", err)
		panic(err)
	}
	log.Print("Catcher stopped")
	err = a.Client.Stop(ctx)
	if err != nil {
		log.Fatal("Error while stopping client", err)
		panic(err)
	}
	log.Print("Client stopped")
}
