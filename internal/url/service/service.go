package service

import (
	"flag"
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

type ConfigVars struct {
	SrvAddr     string
	BaseURL     string
	StoragePath string
}

func SetVars() *ConfigVars {
	serverAddress := flag.String("a", "", "Input server address")
	baseURL := flag.String("b", "", "Input base url")
	fileStoragePath := flag.String("f", "", "Input file storage path")
	flag.Parse()

	if *serverAddress == "" {
		*serverAddress = "localhost:8080"
	}

	if *baseURL == "" {
		*baseURL = "http://localhost:8080"
	}

	if *fileStoragePath == "" {
		*fileStoragePath = ""
	}

	envSrvAddr := os.Getenv("SERVER_ADDRESS")
	if envSrvAddr != "" {
		*serverAddress = envSrvAddr
	}

	envBaseURL := os.Getenv("BASE_URL")
	if envBaseURL != "" {
		*baseURL = envBaseURL
	}

	envFileStorage := os.Getenv("FILE_STORAGE_PATH")
	if envFileStorage != "" {
		*fileStoragePath = envFileStorage
	}

	return &ConfigVars{
		SrvAddr:     *serverAddress,
		BaseURL:     *baseURL,
		StoragePath: *fileStoragePath,
	}
}
