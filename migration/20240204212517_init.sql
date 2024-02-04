-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS public.urls (
    id SERIAL PRIMARY KEY,
    alias_url TEXT UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    visits INTEGER NOT NULL DEFAULT 0,
    count_use INTEGER NOT NULL DEFAULT -1,
    duration TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);
-- +goose Down
DROP TABLE public.urls;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
