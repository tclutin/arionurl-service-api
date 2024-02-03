package app

import (
	"context"
	"github.com/tclutin/ArionURL/internal/config"
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
	router     http.Handler
}

func New(cfg *config.Config, logger *slog.Logger, router http.Handler) *app {
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
