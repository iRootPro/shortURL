package service

import (
	"fmt"
	"os"
  "flag"
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
  if host != "" {
    return host
  }

  baseURLFlag := flag.String("b", "", "Enter base url with port, example: http://localhost:8080")
  flag.Parse()

  if *baseURLFlag != "" {
    return *baseURLFlag
  }
	
  return "http://localhost:8080"
}

func FileStoragePath() string {
  path := os.Getenv("FILE_STORAGE_PATH")
  if path != "" {
    return path
  }

  pathFlag := flag.String("f", "", "Enter path to file storage")
  flag.Parse()

  if *pathFlag != "" {
    return *pathFlag
  }

  return ""

}
