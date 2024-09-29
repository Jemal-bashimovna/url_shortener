package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	urlshortener "shotenedurl"
	"shotenedurl/pkg/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GenerateShortURL() (string, error) {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *Service) CreateURL(originalURL string) (string, error) {

	existURL, err := s.repo.IsExistOriginalURL(originalURL)
	if err != nil {
		return "", err
	}

	if existURL.Id != 0 {
		return existURL.ShortURL, nil
	}

	shortURL, err := s.GenerateShortURL()

	if err != nil {
		return "", err
	}

	url := urlshortener.URL{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
	}

	str, err := s.repo.CreateURL(url)
	if err != nil {
		return "", err
	}
	if str == "" {
		return str, fmt.Errorf("shotUrl is nil")
	}
	return str, nil
}

func (s *Service) GetAll() ([]urlshortener.URL, error) {
	return s.repo.GetAll()
}

func (s *Service) RedirectURL(shortURL string) (string, error) {

	originalURL, err := s.repo.RedirectURL(shortURL)

	if err != nil {
		return "", err
	}
	return originalURL, nil
}

func (s *Service) GetStatsURL(shortURL string) (urlshortener.URLStats, error) {

	return s.repo.GetStatsURL(shortURL)
}
