package service

import (
	"fmt"
	"shotenedurl/models"
	"shotenedurl/pkg/repository"
)

type URLService struct {
	repo repository.URL
}

func NewURLService(repo repository.URL) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) CreateURL(inputURL string) (string, error) {

	existShortURL, err := s.repo.IsExistOriginalURL(inputURL)
	if err != nil {
		return "", err
	}

	if existShortURL != "" {
		return existShortURL, nil
	}

	generatedShortURL, err := generateShortURL()

	if err != nil {
		return "", err
	}

	str, err := s.repo.CreateURL(inputURL, generatedShortURL)
	if err != nil {
		return "", err
	}
	if str == "" {
		return str, fmt.Errorf("shotUrl is nil")
	}
	return str, nil
}

func (s *URLService) RedirectURL(shortURL string) (string, error) {

	originalURL, err := s.repo.GetRedirectURL(shortURL)

	if err != nil {
		return "", err
	}
	return originalURL, nil
}

func (s *URLService) GetStatsURL(shortURL string) (models.URLStats, error) {

	return s.repo.GetStatsURL(shortURL)
}

func (s *URLService) DeleteURL(id int) error {
	return s.repo.DeleteURL(id)
}
