package repository

import (
	"database/sql"
	"fmt"
	"shotenedurl/models"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateURL(inputURL models.InputURL, shortURL string) (string, error) {

	queryURL := fmt.Sprintf("INSERT INTO %s (short_url, original_url) VALUES ($1, $2)", urlsTable)
	_, err := r.db.Exec(queryURL, shortURL, inputURL.OriginalURL)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (r *Repository) IsExistOriginalURL(originalURL models.InputURL) (models.URL, error) {
	var url models.URL
	query := fmt.Sprintf("SELECT * FROM %s WHERE original_url = $1", urlsTable)
	if err := r.db.Get(&url, query, originalURL.OriginalURL); err != nil {
		if err == sql.ErrNoRows {
			return url, nil
		}
		return url, err
	}

	return url, nil
}

func (r *Repository) GetRedirectURL(shortURL string) (string, error) {
	var originalURL string
	var id int

	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf(`
	SELECT id, original_url
	FROM %s WHERE short_url = $1`, urlsTable)
	row := tx.QueryRow(query, shortURL)
	if err := row.Scan(&id, &originalURL); err != nil {
		if err == sql.ErrNoRows {
			tx.Rollback()
			return "", nil
		}
		return "", err
	}

	insertClickQuery := fmt.Sprintf("INSERT INTO %s (url_id) VALUES ($1)", clicksTable)
	_, err = tx.Exec(insertClickQuery, id)

	if err != nil {
		tx.Rollback()
		return "", err
	}

	return originalURL, tx.Commit()
}

func (r *Repository) GetAll() ([]models.URL, error) {
	var urls []models.URL
	query := fmt.Sprintf(`
	SELECT id, short_url, original_url, created_at, expiration_date, deleted_at
	FROM %s `, urlsTable)
	err := r.db.Select(&urls, query)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func (r *Repository) GetStatsURL(shortURL string) (models.URLStats, error) {

	var stats models.URLStats

	query := fmt.Sprintf(`
	SELECT u.short_url, u.original_url, u.created_at, MAX(c.created_at) as last_accessed, COUNT(c.id) as click_count
	FROM %s u 
	LEFT JOIN %s c ON u.id = c.url_id
	WHERE u.short_url = $1 AND u.deleted_at IS NULL
	GROUP BY u.id`, urlsTable, clicksTable)

	if err := r.db.Get(&stats, query, shortURL); err != nil {
		if err == sql.ErrNoRows {
			return stats, fmt.Errorf("short URL not found")
		}
		return stats, err
	}

	return stats, nil
}

func (r *Repository) DeleteURL(id int) error {

	query := fmt.Sprintf(`
    UPDATE %s SET deleted_at=NOW() WHERE id=$1`, urlsTable)
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	updated, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if updated == 0 {
		return fmt.Errorf("url not found with this id")
	}

	return nil
}
