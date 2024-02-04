package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/ArionURL/internal/app"
	"github.com/tclutin/ArionURL/internal/config"
	"github.com/tclutin/ArionURL/internal/controller"
	"github.com/tclutin/ArionURL/internal/repository/shortener/postgres"
	"github.com/tclutin/ArionURL/internal/service/shortener"
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
	client := postgresql.NewClient(context.TODO(), os.Getenv("ARIONURL_DB"))

	//Initializing the router
	router := gin.Default()

	//Initializing the shortener service
	shortenerDBRepo := postgres.NewShortenerRepository(logger, client)
	shortenerDBRepo.InitDB()

	shortenerService := shortener.NewService(logger, cfg, shortenerDBRepo)
	shortenerHandler := controller.NewHandler(logger, cfg, shortenerService)

	shortenerHandler.Register(router)
	app.New(cfg, logger, router).Run()

}
