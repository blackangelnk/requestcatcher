package app

import (
	"context"
	"log"
	"net/http"

	"github.com/blackangelnk/requestcatcher/internal/catcher"
	"github.com/blackangelnk/requestcatcher/internal/client"
	"github.com/blackangelnk/requestcatcher/internal/config"
	"github.com/blackangelnk/requestcatcher/internal/request"
	"github.com/blackangelnk/requestcatcher/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	Client  *client.Client
	Storage storage.Storage
	Catcher *catcher.Catcher
}

func NewApp(cfg *config.Configuration, s storage.Storage) *App {
	app := &App{
		Client:  client.NewClient(cfg, s),
		Storage: s,
		Catcher: catcher.NewCatcher(cfg),
	}
	return app
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
	go func() {
		var cr *request.CaughtRequest
		for {
			cr = <-a.Catcher.Broadcast
			a.Storage.Save(*cr)
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
