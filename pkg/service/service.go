package service

import (
	urlshortener "shotenedurl"
	"shotenedurl/pkg/repository"
)

type URL interface {
	CreateURL(url urlshortener.URL) (string, error)
}

type Service struct {
	URL
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		URL: repo.URL,
	}
}
