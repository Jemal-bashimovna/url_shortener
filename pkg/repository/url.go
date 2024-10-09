package repository

import (
	"context"
	"fmt"
	"urlshortener/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type URLRepository struct {
	db      *pgxpool.Pool
	redisDB *redis.Client
}

func NewURLRepository(db *pgxpool.Pool, redisDB *redis.Client) *URLRepository {
	return &URLRepository{
		db:      db,
		redisDB: redisDB,
	}
}

func (r *URLRepository) CreateURL(inputURL, shortURL string) (string, error) {

	queryURL := fmt.Sprintf("INSERT INTO %s (short_url, original_url) VALUES ($1, $2)", urlsTable)
	_, err := r.db.Exec(context.Background(), queryURL, shortURL, inputURL)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (r *URLRepository) IsExistOriginalURL(inputURL string) (string, error) {
	var shortURL string

	query := fmt.Sprintf("SELECT short_url FROM %s WHERE original_url = $1", urlsTable)
	row := r.db.QueryRow(context.Background(), query, inputURL)

	if err := row.Scan(&shortURL); err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return shortURL, nil
}

func (r *URLRepository) GetRedirectURL(shortURL string) (string, error) {

	// // get shortURL from redis cache
	// cashedURL, err := r.redisDB.Get(ctx, shortURL).Result()
	// if err == nil {
	// 	return cashedURL, err
	// }
	// if err != redis.Nil {
	// 	return "", err
	// }

	var originalURL string
	var id int

	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf(`
	SELECT id, original_url
	FROM %s WHERE short_url = $1`, urlsTable)
	row := tx.QueryRow(context.Background(), query, shortURL)
	if err := row.Scan(&id, &originalURL); err != nil {
		if err == pgx.ErrNoRows {
			tx.Rollback(context.Background())
			return "", nil
		}
		tx.Rollback(context.Background())
		return "", err
	}

	insertClickQuery := fmt.Sprintf("INSERT INTO %s (url_id) VALUES ($1)", clicksTable)
	_, err = tx.Exec(context.Background(), insertClickQuery, id)

	if err != nil {
		tx.Rollback(context.Background())
		return "", err
	}

	// err = r.redisDB.Set(ctx, shortURL, originalURL, 1*time.Hour).Err()
	// if err != nil {
	// 	fmt.Println("Error when saving to Redis cache: ", err)
	// }

	return originalURL, tx.Commit(context.Background())
}

func (r *URLRepository) GetStatsURL(shortURL string) (models.URLStats, error) {

	var stats models.URLStats

	// // Попробуем получить данные из Redis кэша
	// ctx := context.Background()
	// cachedStats, err := r.redisDB.Get(ctx, shortURL).Result()
	// if err == nil {
	// 	// Если данные найдены в кэше, десериализуем их
	// 	if err := json.Unmarshal([]byte(cachedStats), &stats); err != nil {
	// 		return stats, fmt.Errorf("error unmarshaling cached data: %w", err)
	// 	}
	// 	fmt.Println("Данные получены из кэша Redis")
	// 	return stats, nil
	// } else if err != redis.Nil {
	// 	// Если произошла ошибка при обращении к Redis, возвращаем её
	// 	return stats, err
	// }

	query := fmt.Sprintf(`
	SELECT u.short_url, u.original_url, u.created_at, MAX(c.created_at) as last_accessed, COUNT(c.id) as click_count
	FROM %s u 
	LEFT JOIN %s c ON u.id = c.url_id
	WHERE u.short_url = $1 AND u.deleted_at IS NULL
	GROUP BY u.id`, urlsTable, clicksTable)

	row := r.db.QueryRow(context.Background(), query, shortURL)
	if err := row.Scan(&stats.ShortURL, &stats.OriginalURL, &stats.CreatedAt, &stats.LastAccessed, &stats.ClickCount); err != nil {
		if err == pgx.ErrNoRows {
			return stats, fmt.Errorf("short URL not found")
		}
		return stats, err
	}

	// // Сохраняем данные в Redis кэш с истечением времени (например, на 10 минут)
	// cachedData, err := json.Marshal(stats)
	// if err != nil {
	// 	return stats, fmt.Errorf("error marshaling data for Redis: %w", err)
	// }

	// err = r.redisDB.Set(ctx, shortURL, cachedData, 10*time.Minute).Err()
	// if err != nil {
	// 	fmt.Println("Ошибка при сохранении в кэше Redis: ", err)
	// }

	return stats, nil
}

func (r *URLRepository) DeleteURL(id int) error {

	query := fmt.Sprintf(`
    UPDATE %s SET deleted_at=NOW() WHERE id=$1`, urlsTable)
	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	updated := result.RowsAffected()

	if updated == 0 {
		return fmt.Errorf("url not found with this id")
	}

	return nil
}
