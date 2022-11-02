package main

import (
	"log"

	"github.com/irootpro/shorturl/internal/url/handlers"
	"github.com/irootpro/shorturl/internal/url/service"
	"github.com/irootpro/shorturl/internal/url/storage"
)

func InitStorage(cfg *service.ConfigVars) handlers.Storage {
	if cfg.StoragePath == "" {
		return storage.NewStorageMemory()
	}

	if cfg.DSN != "" {
		return storage.NewStorageDB(cfg.DSN)
	}

	storageFile, err := storage.NewStorageFile(cfg.StoragePath)
	if err != nil {
		log.Fatal("create sorage file: ", err)
	}
	return storageFile
}
