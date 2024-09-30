package models

import "time"

type InputURL struct {
	OriginalURL string `json:"original_url" binding:"required"`
}

type URLStats struct {
	ShortURL     string    `json:"short_url" db:"short_url"`
	OriginalURL  string    `json:"original_url" db:"original_url"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	LastAccessed time.Time `json:"last_accessed" db:"last_accessed"`
	ClickCount   int       `json:"click_count" db:"click_count"`
}

type CreateResponse struct {
	ShortURL string
}

type DeleteStatus struct {
	Status string
}
