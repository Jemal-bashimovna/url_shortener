package service

import (
	"urlshortener/models"
	"urlshortener/pkg/repository"
)

type URL interface {
	CreateURL(inputURL string) (string, error)
	RedirectURL(shortURL string) (string, error)
	GetStatsURL(shortURL string) (models.URLStats, error)
	DeleteURL(id int) error
}

type Service struct {
	URL
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		URL: NewURLService(repo.URL),
	}
}
