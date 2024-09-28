package repository

import (
	urlshortener "shotenedurl"

	"github.com/jmoiron/sqlx"
)

type URL interface {
	CreateURL(url urlshortener.URL) (string, error)
}

type Repository struct {
	URL
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		URL: NewUrlPostgres(db),
	}
}
