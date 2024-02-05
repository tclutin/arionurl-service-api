package postgres

import (
	"context"
	"github.com/tclutin/ArionURL/internal/model"
	"github.com/tclutin/ArionURL/pkg/client/postgresql"
	"log/slog"
	"strings"
)

const (
	layer = "shortenerRepository."
)

type shortenerRepository struct {
	client postgresql.Client
	logger *slog.Logger
}

func NewShortenerRepository(logger *slog.Logger, client postgresql.Client) *shortenerRepository {
	logger.Info(layer + "init")
	return &shortenerRepository{
		logger: logger,
		client: client,
	}
}

func (s *shortenerRepository) UpdateShortUrl(ctx context.Context, model *model.URL) error {
	sql := `UPDATE urls SET count_use = $1 WHERE id =  $2`

	s.logger.Info(layer+"updateShortUrl", slog.String("sql", sql))

	_, err := s.client.Exec(ctx, sql, model.Options.CountUse, model.ID)

	if err != nil {
		s.logger.Error(layer+"updateShortUrl", slog.Any("error", err))
		return err
	}

	return nil
}

func (s *shortenerRepository) RemoveUrlByID(ctx context.Context, id uint64) error {
	sql := `DELETE FROM urls WHERE id = $1`

	s.logger.Info(layer+"RemoveUrlByID", slog.String("sql", sql))

	_, err := s.client.Exec(ctx, sql, id)

	if err != nil {
		s.logger.Error(layer+"RemoveUrlByID", slog.Any("error", err))
		return err
	}

	return nil
}

func (s *shortenerRepository) GetUrlByAlias(ctx context.Context, alias string) (*model.URL, error) {
	sql := `SELECT * FROM urls WHERE alias_url = $1`

	s.logger.Info(layer+"GetUrlByAlias", slog.String("sql", sql))

	row := s.client.QueryRow(ctx, sql, alias)

	var url model.URL

	if err := row.Scan(&url.ID, &url.AliasURL, &url.OriginalURL, &url.Options.CountUse, &url.Options.Duration, &url.CreatedAt); err != nil {
		s.logger.Info(layer+"GetUrlByAlias", slog.Any("error", err))
		return nil, err
	}

	return &url, nil
}

func (s *shortenerRepository) CreateAlias(ctx context.Context, model *model.URL) (string, error) {
	sql := `INSERT INTO urls (alias_url, original_url, count_use, duration, created_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING alias_url`

	s.logger.Info(layer+"CreateAlias", slog.String("sql", strings.Replace(sql, "\t", "", -1)))

	row := s.client.QueryRow(ctx, sql, model.AliasURL, model.OriginalURL, model.Options.CountUse, model.Options.Duration, model.CreatedAt)

	var alias string

	if err := row.Scan(&alias); err != nil {
		s.logger.Error(layer+"CreateAlias", slog.Any("error", err))
		return "", err
	}

	return alias, nil
}
