package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/ArionURL/internal/config"
	"github.com/tclutin/ArionURL/internal/controller"
	"github.com/tclutin/ArionURL/internal/repository/postgres"
	"github.com/tclutin/ArionURL/internal/service"
	"github.com/tclutin/ArionURL/pkg/client/postgresql"
	"github.com/tclutin/ArionURL/pkg/logging"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type app struct {
	cfg        *config.Config
	logger     *slog.Logger
	httpServer *http.Server
	router     *gin.Engine
}

func New() *app {
	//Init the cfg
	cfg := config.MustLoad()

	//Init the logger
	logger := logging.InitSlog(cfg.Env)

	//Init the client
	client := postgresql.NewClient(context.Background(), os.Getenv("ARIONURL_DB"))

	//Init the router
	router := gin.Default()

	//Init the shortener layer
	shortenerDBRepo := postgres.NewShortenerRepository(logger, client)
	shortenerService := service.NewShortenerService(logger, cfg, shortenerDBRepo)
	shortenerHandler := controller.NewShortenerHandler(logger, cfg, shortenerService)

	//Register shortener routes
	shortenerHandler.Register(router)

	return &app{
		cfg:    cfg,
		logger: logger,
		router: router,
	}
}

func (a *app) Run() {
	a.startHTTP()
}

func (a *app) startHTTP() {
	a.logger.Info("start HTTP server", slog.String("address", a.cfg.Address))
	a.httpServer = &http.Server{
		Addr:         a.cfg.Address,
		Handler:      a.router,
		WriteTimeout: a.cfg.Timeout,
		ReadTimeout:  a.cfg.Timeout,
	}

	go func() {
		err := a.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			a.logger.Error("the server error has occurred", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	a.logger.Info("shutting down")
	a.shutdownHTTP()
}

func (a *app) shutdownHTTP() {
	if err := a.httpServer.Shutdown(context.Background()); err != nil {
		a.logger.Error("an error occurred during server shutdown", slog.Any("error", err))
		os.Exit(1)
	}
}
