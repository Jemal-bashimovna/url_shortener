package urlshortener

import (
	"errors"
	"time"
)

type URL struct {
	Id             int        `json:"id" db:"id"`
	ShortURL       string     `json:"short_url" bindind:"required" db:"short_url"`
	OriginalURL    string     `json:"original_url" bindind:"required" db:"original_url"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	ExpirationDate *time.Time `json:"expiration_date" db:"expiration_date"`
	DeletedAt      *time.Time `json:"deleted_at" db:"deleted_at"`
}

type InputURL struct {
	OriginalURL string `json:"original_url" bindind:"required"`
	ShortURL    string `json:"short_url" bindind:"required"`
}

func (u *InputURL) ValidateURL(inputURL string) error {
	if len(inputURL) > 250 {
		return errors.New("URL is too long, maximum 250 characters")
	}

	if inputURL[:7] != "http://" && inputURL[:8] != "https://" {
		return errors.New("invalid protocol")
	}
	return nil
}

type URLStats struct {
	ShortURL     string     `json:"short_url"`
	OriginalURL  string     `json:"original_url"`
	CreatedAt    time.Time  `json:"created_at"`
	LastAccessed *time.Time `json:"last_accessed"`
	ClickCount   int        `json:"click_count"`
}
