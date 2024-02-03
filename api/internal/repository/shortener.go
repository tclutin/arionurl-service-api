package repository

import (
	"context"
	"github.com/tclutin/ArionURL/internal/domain/shortener"
	"github.com/tclutin/ArionURL/pkg/client/postgresql"
	"log"
	"log/slog"
)

type shortenerRepository struct {
	client postgresql.Client
	logger *slog.Logger
}

func (s *shortenerRepository) RemoveShortUrl(alias string) error {
	//TODO implement me
	panic("implement me")
}

func (s *shortenerRepository) UpdateShortUrl(dto shortener.URL) {
	//TODO implement me
	panic("implement me")
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

func (s *shortenerRepository) GetUrlByAlias(alias string) (shortener.URL, error) {
	sql := `SELECT * FROM urls WHERE alias_url = $1`

	row := s.client.QueryRow(context.Background(), sql, alias)

	var url shortener.URL

	if err := row.Scan(&url.ID, &url.UserID, &url.AliasURL, &url.OriginalURL, &url.Options.Visits, &url.Options.CountUse, &url.Options.Duration, &url.CreatedAt); err != nil {
		return url, err
	}

	return url, nil
}

func (s *shortenerRepository) CreateAlias(model shortener.URL) (string, error) {
	sql := `INSERT INTO urls (alias_url, original_url, visits, count_use, duration, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING alias_url`

	row := s.client.QueryRow(context.Background(), sql, model.AliasURL, model.OriginalURL, model.Options.Visits, model.Options.CountUse, model.Options.Duration, model.CreatedAt)

	var alias string

	if err := row.Scan(&alias); err != nil {
		log.Fatalln(err)
	}

	return alias, nil
}
