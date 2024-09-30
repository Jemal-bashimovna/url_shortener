package repository

import (
	"shotenedurl/models"

	"github.com/jmoiron/sqlx"
)

type URL interface {
	CreateURL(inputURL, generatedShortURL string) (string, error)
	IsExistOriginalURL(inputURL string) (string, error)
	GetRedirectURL(shortURL string) (string, error)
	GetStatsURL(shortURL string) (models.URLStats, error)
	DeleteURL(id int) error
}

type Repository struct {
	URL
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		URL: NewURLRepository(db),
	}
}
