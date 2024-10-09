package repository

import (
	"urlshortener/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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

func NewRepository(db *pgxpool.Pool, redisDB *redis.Client) *Repository {
	return &Repository{
		URL: NewURLRepository(db, redisDB),
	}
}
