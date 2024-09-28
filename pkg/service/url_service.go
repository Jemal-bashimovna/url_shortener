package service

import (
	urlshortener "shotenedurl"
	"shotenedurl/pkg/repository"
)

type UrlService struct {
	repo repository.UrlPostgres
}

func NewUrlService(repo repository.UrlPostgres) *UrlService {
	return &UrlService{repo: repo}
}

func (s *UrlService) CreateURL(url urlshortener.URL) (string, error) {
	return s.repo.CreateURL(url)
}
