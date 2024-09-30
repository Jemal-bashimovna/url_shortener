package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"shotenedurl/models"
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

func (s *Service) CreateURL(inputURL models.InputURL) (string, error) {

	existURL, err := s.repo.IsExistOriginalURL(inputURL)
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

	str, err := s.repo.CreateURL(inputURL, shortURL)
	if err != nil {
		return "", err
	}
	if str == "" {
		return str, fmt.Errorf("shotUrl is nil")
	}
	return str, nil
}

func (s *Service) GetAll() ([]models.URL, error) {
	return s.repo.GetAll()
}

func (s *Service) RedirectURL(shortURL string) (string, error) {

	originalURL, err := s.repo.GetRedirectURL(shortURL)

	if err != nil {
		return "", err
	}
	return originalURL, nil
}

func (s *Service) GetStatsURL(shortURL string) (models.URLStats, error) {

	return s.repo.GetStatsURL(shortURL)
}

func (s *Service) DeleteURL(id int) error {
	return s.repo.DeleteURL(id)
}
