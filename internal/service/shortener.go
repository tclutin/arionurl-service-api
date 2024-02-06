package service

import (
	"context"
	"errors"
	"github.com/tclutin/ArionURL/pkg/utils"
	"github.com/tclutin/arionurl-service-api/internal/config"
	"github.com/tclutin/arionurl-service-api/internal/controller/dto"
	"github.com/tclutin/arionurl-service-api/internal/model"
	"log/slog"
	"net/url"
	"strings"
	"time"
)

const (
	layer = "shortenerService."
)

type shortenerRepository interface {
	CreateAlias(ctx context.Context, model *model.URL) (string, error)
	GetUrlByAlias(ctx context.Context, alias string) (*model.URL, error)
	RemoveUrlByID(ctx context.Context, id uint64) error
	UpdateShortUrl(ctx context.Context, model *model.URL) error
}

type shortenerService struct {
	logger *slog.Logger
	cfg    *config.Config
	repo   shortenerRepository
}

func NewShortenerService(logger *slog.Logger, cfg *config.Config, repo shortenerRepository) *shortenerService {
	logger.Info(layer + "init")
	return &shortenerService{logger: logger, cfg: cfg, repo: repo}
}

func (s *shortenerService) LookShortUrl(ctx context.Context, alias string) (*model.URL, error) {
	s.logger.Info(layer + "LookShortUrl")
	url, err := s.repo.GetUrlByAlias(ctx, alias)
	if err != nil {
		return nil, errors.New("alias not found")
	}

	if url.Options.Duration.Before(time.Now()) {

		err = s.repo.RemoveUrlByID(ctx, url.ID)
		if err != nil {
			return nil, errors.New("deletion error")
		}

		return nil, errors.New("url expired")
	}

	if url.Options.CountUse == 0 {

		err = s.repo.RemoveUrlByID(ctx, url.ID)
		if err != nil {
			return nil, errors.New("deletion error")
		}
		return nil, errors.New("count is poor")
	}

	if url.Options.CountUse > 0 {
		url.Options.CountUse--
	}

	err = s.repo.UpdateShortUrl(ctx, url)
	if err != nil {
		return nil, errors.New("failed to update short url")
	}

	return url, nil
}

func (s *shortenerService) CreateShortUrl(ctx context.Context, dto dto.CreateUrlRequest) (string, error) {
	s.logger.Info(layer + "CreateShortUrl")
	dto.Duration = strings.ReplaceAll(dto.Duration, "-", "")

	if !s.validateOriginalURL(dto.OriginalURL) {
		return "", errors.New("invalid format of the original url")
	}

	duration, err := time.ParseDuration(dto.Duration)
	if err != nil {
		return "", errors.New("invalid format of the time")
	}

	if !s.validateDuration(duration) {
		return "", errors.New("count of hours >= 720h")
	}

	if dto.CountUse <= 0 {
		dto.CountUse = -1
	}

	currentTime := time.Now().UTC()
	expirationTime := currentTime.Add(duration)

	options := model.URLOptions{
		Duration: expirationTime,
		CountUse: dto.CountUse,
	}

	url := &model.URL{
		AliasURL:    s.generateAlias(s.cfg.SizeShortUrl),
		OriginalURL: dto.OriginalURL,
		Options:     options,
		CreatedAt:   currentTime,
	}

	alias, err := s.repo.CreateAlias(ctx, url)
	if err != nil {
		return "", errors.New("alias creation error")
	}

	return alias, nil
}

func (s *shortenerService) validateDuration(duration time.Duration) bool {
	s.logger.Info(layer + "validateDuration")
	maxDuration := 720 * time.Hour
	if duration < maxDuration {
		return true
	}
	return false
}

func (s *shortenerService) validateOriginalURL(originalURL string) bool {
	s.logger.Info(layer + "validateOriginalURL")
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return false
	}
	return true
}

func (s *shortenerService) generateAlias(size int64) string {
	s.logger.Info(layer + "generateAlias")
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	alias := make([]rune, size)
	for i := range alias {
		alias[i] = chars[utils.NewCryptoRand(int64(len(chars)))]
	}
	return string(alias)
}
