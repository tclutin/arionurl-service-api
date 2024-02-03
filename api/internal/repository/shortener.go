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
    		user_id INTEGER UNIQUE,
    		alias_url TEXT NOT NULL,
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

func (s *shortenerRepository) GetByAlias(alias string) (*shortener.URL, error) {
	//TODO implement me
	panic("implement me")
}

func (s *shortenerRepository) CreateAlias(dto shortener.CreateUrlDTO) (string, error) {
	//TODO implement me
	panic("implement me")
}
