package repository

import (
	"shotenedurl/models"

	"github.com/jackc/pgx/v5/pgxpool"
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

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		URL: NewURLRepository(db),
	}
}
