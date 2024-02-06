package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/arionurl-service-api/internal/config"
	"github.com/tclutin/arionurl-service-api/internal/controller/dto"
	"github.com/tclutin/arionurl-service-api/internal/controller/middleware"
	"github.com/tclutin/arionurl-service-api/internal/model"
	"log/slog"
	"net/http"
)

const (
	layer              = "shortenerHandler."
	createAliasURL     = "/aliases"
	redirectToAliasURL = "/:alias"
)

type shortenerService interface {
	CreateShortUrl(ctx context.Context, dto dto.CreateUrlRequest) (string, error)
	LookShortUrl(ctx context.Context, alias string) (*model.URL, error)
}

type shortenerHandler struct {
	cfg     *config.Config
	logger  *slog.Logger
	service shortenerService
}

func NewShortenerHandler(logger *slog.Logger, cfg *config.Config, service shortenerService) *shortenerHandler {
	logger.Info(layer + "init")
	return &shortenerHandler{logger: logger, cfg: cfg, service: service}
}

func (h *shortenerHandler) Register(router *gin.Engine) {
	router.Use(middleware.RateLimiter())
	router.POST(createAliasURL, h.CreateAlias)
	router.GET(redirectToAliasURL, h.RedirectToAlias)
}

func (h *shortenerHandler) CreateAlias(c *gin.Context) {
	h.logger.Info(layer + "CreateAlias")
	var dto dto.CreateUrlRequest
	if err := c.BindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	alias, err := h.service.CreateShortUrl(context.Background(), dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fullUrl := fmt.Sprintf("http://%s/%s", h.cfg.Address, alias)
	c.JSON(http.StatusCreated, gin.H{"alias": fullUrl})
	return
}

func (h *shortenerHandler) RedirectToAlias(c *gin.Context) {
	h.logger.Info(layer + "RedirectToAlias")
	alias := c.Param("alias")

	url, err := h.service.LookShortUrl(context.Background(), alias)
	if err != nil {
		c.Redirect(http.StatusFound, h.cfg.BaseRedirect)
		return
	}

	c.Redirect(http.StatusFound, url.OriginalURL)
	return
}
