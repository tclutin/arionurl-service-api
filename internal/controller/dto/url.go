package dto

type CreateUrlRequest struct {
	OriginalURL string `json:"original_url" binding:"required"`
	Duration    string `json:"duration" binding:"required"`
	CountUse    int    `json:"count_use"`
}
