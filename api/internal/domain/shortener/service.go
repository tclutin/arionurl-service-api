package shortener

import (
	"errors"
	"github.com/tclutin/ArionURL/pkg/utils"
	"log/slog"
	"net/url"
	"strings"
	"time"
)

type Repository interface {
	CreateAlias(model URL) (string, error)
	GetUrlByAlias(alias string) (URL, error)
	RemoveShortUrl(alias string) error
	UpdateShortUrl(model URL)
}

type service struct {
	logger *slog.Logger
	repo   Repository
}

func NewService(logger *slog.Logger, repo Repository) *service {
	return &service{logger: logger, repo: repo}
}

func (s *service) CreateShortUrl(dto CreateUrlDTO) (string, error) {
	if !s.validateOriginalURL(dto.OriginalURL) {
		return "", errors.New("invalid format of the original url")
	}

	dto.Duration = strings.ReplaceAll(dto.Duration, "-", "")

	duration, err := time.ParseDuration(dto.Duration)
	if err != nil {
		return "", errors.New("invalid format of the time")
	}

	if !s.validateDuration(duration) {
		return "", errors.New("count of hours >= 720h")
	}

	currentTime := time.Now().UTC()
	expirationTime := currentTime.Add(duration)

	options := URLOptions{
		Duration: expirationTime,
		CountUse: dto.CountUse,
	}

	url := URL{
		AliasURL:    s.generateAlias(6),
		OriginalURL: dto.OriginalURL,
		Options:     options,
		CreatedAt:   currentTime,
	}

	alias, err := s.repo.CreateAlias(url)
	if err != nil {
		return "", errors.New("alias creation error")
	}

	return alias, nil
}

func (s *service) LookShortUrl(alias string) (URL, error) {
	url, err := s.repo.GetUrlByAlias(alias)
	if err != nil {
		return URL{}, errors.New("alias not found")
	}

	if url.Options.Duration.Before(time.Now()) {
		return URL{}, errors.New("URL expired")
	}

	return url, nil
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
