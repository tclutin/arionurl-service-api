package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tclutin/ArionURL/internal/config"
	"github.com/tclutin/ArionURL/internal/domain/shortener"
	"log/slog"
	"net/http"
)

const (
	createAliasURL  = "/aliases"
	redirectToAlias = "/:alias"
)

type ShortenerService interface {
	CreateShortUrl(dto shortener.CreateUrlDTO) (string, error)
	LookShortUrl(alias string) (shortener.URL, error)
}

type handler struct {
	cfg     *config.Config
	logger  *slog.Logger
	service ShortenerService
}

func NewHandler(logger *slog.Logger, cfg *config.Config, service ShortenerService) *handler {
	return &handler{logger: logger, cfg: cfg, service: service}
}

func (h *handler) Register(router *gin.Engine) {
	router.POST(createAliasURL, h.CreateAlias)
	router.GET(redirectToAlias, h.RedirectToAlias)
}

func (h *handler) CreateAlias(c *gin.Context) {
	var dto shortener.CreateUrlDTO
	if err := c.BindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrl, err := h.service.CreateShortUrl(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"alias": shortUrl})
	return
}

func (h *handler) RedirectToAlias(c *gin.Context) {
	alias := c.Param("alias")
	url, err := h.service.LookShortUrl(alias)

	if err != nil {
		c.Redirect(http.StatusFound, h.cfg.BaseRedirect)
		return
	}
	c.Redirect(http.StatusFound, url.OriginalURL)
	return
}
