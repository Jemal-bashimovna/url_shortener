package repository

import (
	"database/sql"
	"fmt"
	"shotenedurl/models"

	"github.com/jmoiron/sqlx"
)

type URLRepository struct {
	db *sqlx.DB
}

func NewURLRepository(db *sqlx.DB) *URLRepository {
	return &URLRepository{db: db}
}

func (r *URLRepository) CreateURL(inputURL, shortURL string) (string, error) {

	queryURL := fmt.Sprintf("INSERT INTO %s (short_url, original_url) VALUES ($1, $2)", urlsTable)
	_, err := r.db.Exec(queryURL, shortURL, inputURL)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (r *URLRepository) IsExistOriginalURL(inputURL string) (string, error) {
	var shortURL string
	query := fmt.Sprintf("SELECT short_url FROM %s WHERE original_url = $1", urlsTable)
	row := r.db.QueryRow(query, inputURL)

	if err := row.Scan(&shortURL); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return shortURL, nil
}

func (r *URLRepository) GetRedirectURL(shortURL string) (string, error) {
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
		tx.Rollback()
		return "", err
	}

	insertClickQuery := fmt.Sprintf("INSERT INTO %s (url_id) VALUES ($1)", clicksTable)
	fmt.Println(insertClickQuery)
	_, err = tx.Exec(insertClickQuery, id)

	if err != nil {
		tx.Rollback()
		return "", err
	}

	return originalURL, tx.Commit()
}

func (r *URLRepository) GetStatsURL(shortURL string) (models.URLStats, error) {

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

func (r *URLRepository) DeleteURL(id int) error {

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
