package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/ArionURL/internal/config"
	"github.com/tclutin/ArionURL/internal/service/shortener"
	"log/slog"
	"net/http"
)

const (
	layer              = "shortenerHandler."
	createAliasURL     = "/aliases"
	redirectToAliasURL = "/:alias"
)

type ShortenerService interface {
	CreateShortUrl(ctx context.Context, dto shortener.CreateUrlDTO) (string, error)
	LookShortUrl(ctx context.Context, alias string) (*shortener.URL, error)
}

type handler struct {
	cfg     *config.Config
	logger  *slog.Logger
	service ShortenerService
}

func NewHandler(logger *slog.Logger, cfg *config.Config, service ShortenerService) *handler {
	logger.Info(layer + "init")
	return &handler{logger: logger, cfg: cfg, service: service}
}

func (h *handler) Register(router *gin.Engine) {
	router.POST(createAliasURL, h.CreateAlias)
	router.GET(redirectToAliasURL, h.RedirectToAlias)
}

func (h *handler) CreateAlias(c *gin.Context) {
	var dto shortener.CreateUrlDTO
	if err := c.BindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrl, err := h.service.CreateShortUrl(context.Background(), dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"alias": shortUrl})
	return
}

func (h *handler) RedirectToAlias(c *gin.Context) {
	alias := c.Param("alias")
	url, err := h.service.LookShortUrl(context.Background(), alias)

	if err != nil {
		c.Redirect(http.StatusFound, h.cfg.BaseRedirect)
		return
	}
	c.Redirect(http.StatusFound, url.OriginalURL)
	return
}
