package user

import "time"

type User struct {
	ID         uint64    `json:"id"`
	Username   string    `json:"username"`
	TelegramID string    `json:"telegram_id"`
	CreatedAt  time.Time `json:"created_at"`
}
