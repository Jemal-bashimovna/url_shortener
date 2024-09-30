package handler

import (
	"fmt"
	"net/url"
)

func validateURL(inputURL string) error {
	// if len(inputURL) > 250 {
	// 	return errors.New("URL is too long, maximum 250 characters")
	// }

	// if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
	// 	return errors.New("invalid protocol")
	// }
	// return nil

	parsedURL, err := url.ParseRequestURI(inputURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %s", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("invalid URL protocol: %s", parsedURL.Scheme)
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("host is missing in the URL")
	}

	return nil
}
