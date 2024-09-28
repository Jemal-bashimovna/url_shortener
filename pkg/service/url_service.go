package services

import (
	"crypto/rand"
	"encoding/hex"
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

	existURL, err := s.repo.IsExistURL(originalURL)

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

	return str, nil
}

func (s *Service) GetAll() ([]urlshortener.URL, error) {
	return s.repo.GetAll()
}
