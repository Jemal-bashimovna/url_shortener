package repository

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	urlshortener "shotenedurl"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type UrlPostgres struct {
	db *sqlx.DB
}

func NewUrlPostgres(db *sqlx.DB) *UrlPostgres {
	return &UrlPostgres{db: db}
}

func (r *UrlPostgres) CreateURL(url urlshortener.URL) (string, error) {
	str, err := GenerateShortURL()
	if err != nil {
		return "", err
	}
	shortURL := viper.GetString("domain") + str
	var id int
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}
	fmt.Println(shortURL)
	queryURL := fmt.Sprintf("INSERT INTO %s (short_url, original_url) VALUES ($1, $2) RETURNING id", urlsTable)
	row := tx.QueryRow(queryURL, shortURL, url.OriginalURL)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return "", err
	}

	queryClick := fmt.Sprintf("INSERT INTO %s (url_id) VALUES ($1) ", clicksTable)
	_, err = tx.Exec(queryClick, id)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	return shortURL, tx.Commit()
}

func GenerateShortURL() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
