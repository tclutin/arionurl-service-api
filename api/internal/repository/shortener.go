package repository

import (
	"context"
	"github.com/tclutin/ArionURL/internal/service/shortener"
	"github.com/tclutin/ArionURL/pkg/client/postgresql"
	"log"
	"log/slog"
)

type shortenerRepository struct {
	client postgresql.Client
	logger *slog.Logger
}

func NewShortenerRepo(logger *slog.Logger, client postgresql.Client) *shortenerRepository {
	return &shortenerRepository{
		logger: logger,
		client: client,
	}
}

func (s *shortenerRepository) InitDB() {
	users := `CREATE TABLE IF NOT EXISTS public.users (
    		id SERIAL PRIMARY KEY,
    		username TEXT NOT NULL,
    		telegram_id TEXT,
    		created_at TIMESTAMP NOT NULL 	
		)`

	urls := `CREATE TABLE IF NOT EXISTS public.urls (
    		id SERIAL PRIMARY KEY,
    		user_id INTEGER,
    		alias_url TEXT UNIQUE NOT NULL,
    		original_url TEXT NOT NULL, 
    		visits INTEGER NOT NULL DEFAULT 0,
    		count_use INTEGER NOT NULL DEFAULT -1,
    		duration TIMESTAMP NOT NULL,
    	    created_at TIMESTAMP NOT NULL,
    	    FOREIGN KEY (user_id) REFERENCES public.users(id)
			)`

	_, err := s.client.Exec(context.Background(), users)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = s.client.Exec(context.Background(), urls)
	if err != nil {
		log.Fatalln(err)
	}
}

func (s *shortenerRepository) UpdateShortUrl(ctx context.Context, entity *shortener.URL) error {
	sql := `UPDATE urls SET visits = $1, count_use = $2 WHERE id =  $3`

	_, err := s.client.Exec(ctx, sql, entity.Options.Visits, entity.Options.CountUse, entity.ID)

	if err != nil {
		return err
	}
	return nil
}

func (s *shortenerRepository) RemoveUrlByID(ctx context.Context, id uint64) error {
	sql := `DELETE FROM urls WHERE id = $1`
	_, err := s.client.Exec(ctx, sql, id)

	if err != nil {
		return err
	}
	return nil
}

func (s *shortenerRepository) GetUrlByAlias(ctx context.Context, alias string) (*shortener.URL, error) {
	sql := `SELECT * FROM urls WHERE alias_url = $1`

	row := s.client.QueryRow(ctx, sql, alias)

	var url shortener.URL

	if err := row.Scan(&url.ID, &url.UserID, &url.AliasURL, &url.OriginalURL, &url.Options.Visits, &url.Options.CountUse, &url.Options.Duration, &url.CreatedAt); err != nil {
		return nil, err
	}
	return &url, nil
}

func (s *shortenerRepository) CreateAlias(ctx context.Context, model *shortener.URL) (string, error) {
	sql := `INSERT INTO urls (alias_url, original_url, visits, count_use, duration, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING alias_url`

	row := s.client.QueryRow(ctx, sql, model.AliasURL, model.OriginalURL, model.Options.Visits, model.Options.CountUse, model.Options.Duration, model.CreatedAt)

	var alias string

	if err := row.Scan(&alias); err != nil {
		log.Fatalln(err)
	}
	return alias, nil
}
