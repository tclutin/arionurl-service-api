package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/ArionURL/internal/app"
	"github.com/tclutin/ArionURL/internal/config"
	"github.com/tclutin/ArionURL/internal/controller"
	"github.com/tclutin/ArionURL/internal/domain/shortener"
	"github.com/tclutin/ArionURL/internal/repository"
	"github.com/tclutin/ArionURL/pkg/client/postgresql"
	"github.com/tclutin/ArionURL/pkg/logging"
	"os"
)

func main() {
	//Initializing the config
	cfg := config.MustLoad()

	//Initializing the logger
	logger := logging.InitSlog(cfg.Env)

	//Initializing the pgxpool
	pgxPool := postgresql.NewClient(context.Background(), os.Getenv("ARIONURL_DB"))

	//Initializing the router
	router := gin.Default()

	//Initializing the shortener service
	shortenerRepo := repository.NewShortenerRepo(logger, pgxPool)
	shortenerRepo.InitDB()
	shortenerService := shortener.NewService(logger, shortenerRepo)
	shortenerHandler := controller.NewHandler(logger, shortenerService)

	shortenerHandler.Register(router)
	app.New(cfg, logger, router).Run()

}
