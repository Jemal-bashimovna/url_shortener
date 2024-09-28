package urlshortener

import "errors"

type URL struct {
	Id             int    `json:"id"`
	ShortURL       string `json:"short_url" bindind:"required"`
	OriginalURL    string `json:"original_url" bindind:"required"`
	CreatedAt      string `json:"created_at"`
	ExpirationDate string `json:"exp_date"`
	DeletedAt      string `json:"deleted_at"`
}

func (u *URL) ValidateURL(inputURL string) error {
	if len(inputURL) > 250 {
		return errors.New("URL is too long, maximum 250 characters")
	}

	if inputURL[:7] != "http://" && inputURL[:8] != "https://" {
		return errors.New("invalid protocol")
	}
	return nil
}
