package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/ArionURL/internal/app"
	"github.com/tclutin/ArionURL/internal/config"
	"github.com/tclutin/ArionURL/internal/controller"
	"github.com/tclutin/ArionURL/internal/repository/postgres"
	"github.com/tclutin/ArionURL/internal/service/shortener"
	"github.com/tclutin/ArionURL/pkg/client/postgresql"
	"github.com/tclutin/ArionURL/pkg/logging"
	"net/url"
	"os"
	"unsafe"
)

func main() {
	fmt.Println(unsafe.Sizeof(url.URL{}))
	//Initializing the config
	cfg := config.MustLoad()

	//Initializing the logger
	logger := logging.InitSlog(cfg.Env)

	//Initializing the pgxpool
	pgxPool := postgresql.NewClient(context.TODO(), os.Getenv("ARIONURL_DB"))

	//Initializing the router
	router := gin.Default()

	//Initializing the shortener service
	shortenerRepo := postgres.NewShortenerRepo(logger, pgxPool)
	shortenerRepo.InitDB()

	shortenerService := shortener.NewService(logger, cfg, shortenerRepo)
	shortenerHandler := controller.NewHandler(logger, cfg, shortenerService)

	shortenerHandler.Register(router)
	app.New(cfg, logger, router).Run()

}
