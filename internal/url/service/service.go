package service

import (
	"fmt"
	"os"

	"github.com/irootpro/shorturl/internal/url/storage"
)

func ShortURLByID(links []storage.LinkEntity, id string) (string, error) {
	for _, v := range links {
		if v.ID == id {
			return v.OriginalURL, nil
		}
	}

	return "", fmt.Errorf("link not found")
}

func BaseURL() string {
	host := os.Getenv("BASE_URL")
	if host == "" {
		host = "http://localhost:8080"
	}

	return host
}

func FileStorageEnv() string {
	return os.Getenv("FILE_STORAGE_PATH")
}
