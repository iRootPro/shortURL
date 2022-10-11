package service

import (
	"fmt"

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
