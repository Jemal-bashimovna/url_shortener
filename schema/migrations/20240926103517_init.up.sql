CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_url VARCHAR(255) UNIQUE,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expiration_date TIMESTAMP,
    deleted_at TIMESTAMP

);

CREATE TABLE clicks(
    id SERIAL PRIMARY KEY,
    url_id INT REFERENCES urls(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);