package shortener

import (
	"context"
	"errors"
	"github.com/tclutin/ArionURL/internal/config"
	"github.com/tclutin/ArionURL/pkg/utils"
	"log/slog"
	"net/url"
	"strings"
	"time"
)

type Repository interface {
	CreateAlias(ctx context.Context, model *URL) (string, error)
	GetUrlByAlias(ctx context.Context, alias string) (*URL, error)
	RemoveUrlByID(ctx context.Context, id uint64) error
	UpdateShortUrl(ctx context.Context, model *URL) error
}

type service struct {
	logger *slog.Logger
	cfg    *config.Config
	repo   Repository
}

func NewService(logger *slog.Logger, cfg *config.Config, repo Repository) *service {
	return &service{logger: logger, cfg: cfg, repo: repo}
}

func (s *service) LookShortUrl(ctx context.Context, alias string) (*URL, error) {
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

	url.Options.Visits++

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

func (s *service) CreateShortUrl(ctx context.Context, dto CreateUrlDTO) (string, error) {
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

	options := URLOptions{
		Duration: expirationTime,
		CountUse: dto.CountUse,
	}

	url := &URL{
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

func (s *service) validateDuration(duration time.Duration) bool {
	maxDuration := 720 * time.Hour
	if duration < maxDuration {
		return true
	}
	return false
}

func (s *service) validateOriginalURL(originalURL string) bool {
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return false
	}
	return true
}

func (s *service) generateAlias(size int64) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	alias := make([]rune, size)
	for i := range alias {
		alias[i] = chars[utils.NewCryptoRand(int64(len(chars)))]
	}
	return string(alias)
}
