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

func (r *Repository) IsExistOriginalURL(originalURL string) (urlshortener.URL, error) {
	var url urlshortener.URL
	query := fmt.Sprintf("SELECT * FROM %s WHERE original_url = $1", urlsTable)
	if err := r.db.Get(&url, query, originalURL); err != nil {
		if err == sql.ErrNoRows {
			return url, nil
		}
		return url, err
	}

	return url, nil
}

func (r *Repository) RedirectURL(shortURL string) (string, error) {
	var originalURL string
	var id int

	query := fmt.Sprintf(`
	SELECT id, original_url
	FROM %s WHERE short_url = $1`, urlsTable)
	row := r.db.QueryRow(query, shortURL)
	if err := row.Scan(&id, &originalURL); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	insertClickQuery := fmt.Sprintf("INSERT INTO %s (url_id) VALUES ($1)", clicksTable)
	_, err := r.db.Exec(insertClickQuery, id)

	if err != nil {
		return "", err
	}

	return originalURL, nil
}

func (r *Repository) IsExistShortURL(shortURL string) (urlshortener.URL, error) {
	var url urlshortener.URL
	query := fmt.Sprintf(`
	SELECT ut.id, ut.short_url, ut.original_url, ut.created_at, ut.expiration_date, ut.deleted_at 
	FROM %s ut WHERE ut.short_url = $1`, urlsTable)
	err := r.db.Get(&url, query, shortURL)
	if err != nil && err == sql.ErrNoRows {
		return url, nil
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

func (r *Repository) GetStatsURL(shortURL string) (urlshortener.URLStats, error) {

	var stats urlshortener.URLStats

	query := fmt.Sprintf(`
	SELECT u.short_url, u.original_url, u.created_at, MAX(c.created_at) as last_accessed, COUNT(c.id) as click_count
	FROM %s u 
	LEFT JOIN %s c ON u.id = c.url_id
	WHERE u.short_url = $1
	GROUP BY u.id`, urlsTable, clicksTable)

	if err := r.db.Get(&stats, query, shortURL); err != nil {
		if err == sql.ErrNoRows {
			return stats, fmt.Errorf("short URL not found")
		}
		return stats, err
	}

	return stats, nil
}
