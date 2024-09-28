package repository

import (
	"database/sql"
	"fmt"
	urlshortener "shotenedurl"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateURL(url urlshortener.URL) (string, error) {

	var id int
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	queryURL := fmt.Sprintf("INSERT INTO %s (short_url, original_url) VALUES ($1, $2) RETURNING id", urlsTable)
	row := tx.QueryRow(queryURL, url.ShortURL, url.OriginalURL)
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

	return url.ShortURL, tx.Commit()
}

func (r *Repository) IsExistURL(originalURL string) (urlshortener.URL, error) {
	var url urlshortener.URL
	query := fmt.Sprintf("SELECT * FROM %s WHERE original_url = $1", urlsTable)
	row := r.db.QueryRow(query, originalURL)
	err := row.Scan(&url.Id, &url.ShortURL, &url.OriginalURL, &url.CreatedAt, &url.ExpirationDate, &url.DeletedAt)
	if err != nil {
		return url, err
	}
	return url, nil
}

func (r *Repository) IsExistShortURL(shortURL string) (urlshortener.URL, error) {
	var url urlshortener.URL
	query := fmt.Sprintf(`
	SELECT ut.id, ut.short_url, ut.original_url, ut.created_at, ut.expiration_date, ut.deleted_at 
	FROM %s ut WHERE ut.short_url = $1`, urlsTable)
	err := r.db.Get(&url, query, shortURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return url, nil
		}
	}
	return url, err
}

func (r *Repository) GetAll() ([]urlshortener.URL, error) {
	var url []urlshortener.URL
	query := fmt.Sprintf(`
	SELECT ut.id, ut.short_url, ut.original_url, ut.created_at, ut.expiration_date, ut.deleted_at
	FROM %s ut `, urlsTable)
	err := r.db.Select(&url, query)
	if err != nil {
		return nil, err
	}
	return url, nil
}
