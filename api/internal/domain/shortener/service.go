package shortener

import (
	"errors"
	"log/slog"
	"net/url"
	"time"
)

type Repository interface {
	CreateAlias(dto CreateUrlDTO) (string, error)
	GetByAlias(alias string) (*URL, error)
}

type service struct {
	logger *slog.Logger
	repo   Repository
}

func NewService(logger *slog.Logger, repo Repository) *service {
	return &service{logger: logger, repo: repo}
}

func (s *service) CreateURL(dto CreateUrlDTO) (string, error) {
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
	return "", nil
}

func (s *service) GetURLByAlias(alias string) (*URL, error) {
	return nil, nil
}

func (s *service) validateDuration(duration time.Duration) bool {
	maxDuration := 720 * time.Hour
	if duration < maxDuration {
		return true
	}
	return false
}

func (s *service) validateOriginalURL(original string) bool {
	_, err := url.ParseRequestURI(original)
	if err != nil {
		return false
	}
	return true
}
