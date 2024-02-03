package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/tclutin/ArionURL/internal/app"
	"github.com/tclutin/ArionURL/internal/config"
	"github.com/tclutin/ArionURL/pkg/client/postgresql"
	"github.com/tclutin/ArionURL/pkg/logging"
	"os"
)

func main() {
	//Initializing the config
	cfg := config.MustLoad()

	//Initializing the logger
	logger := logging.InitSlog(cfg.Env)

	//Initial the pgxpool
	pgxPool := postgresql.NewClient(context.Background(), os.Getenv("ARIONURL_DB"))
	fmt.Println(pgxPool)

	//Initial the router
	router := chi.NewRouter()

	//Initial the app
	app.New(cfg, logger, router).Run()
}
