package main

import (
	"github.com/irootpro/shorturl/internal/url/handlers"
	"github.com/irootpro/shorturl/internal/url/service"
	"github.com/irootpro/shorturl/internal/url/storage"
)

func InitStorage(cfg *service.ConfigVars) handlers.Storage {
	if cfg.StoragePath == "" {
		return storage.NewStorageMemory()
	}

	return storage.NewStorageFile(cfg.StoragePath)
}
