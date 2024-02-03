package shortener

import "time"

type URL struct {
	ID          uint64     `json:"id"`
	UserID      uint64     `json:"user_id"`
	AliasURL    string     `json:"alias_url"`
	OriginalURL string     `json:"original_url"`
	Options     URLOptions `json:"options"`
	CreatedAt   time.Time  `json:"created_at"`
}

type URLOptions struct {
	Visits   int           `json:"visits"`
	Duration time.Duration `json:"duration"`
}
