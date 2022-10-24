package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/irootpro/shorturl/internal/url/handlers"
	"github.com/irootpro/shorturl/internal/url/service"
)

func main() {
	cfg := service.SetVars()
	storage := InitStorage(cfg)
	serverHandler := handlers.NewServerHandler(cfg, storage)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())
	e.Use(middleware.Decompress())
	e.GET("/:hash", serverHandler.GetURL)
	e.POST("/", serverHandler.PostURL)
	e.POST("/api/shorten", serverHandler.PostURLJSON)
	if err := e.Start(cfg.SrvAddr); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
