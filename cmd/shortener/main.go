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
	getURLHandler := handlers.NewGetURLHandler(cfg)
	postURLHandler := handlers.NewPostURLHandler(cfg)
	postURLJSONHandler := handlers.NewPostURLJSONHandler(cfg)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())
	e.GET("/:hash", getURLHandler.GetURL)
	e.POST("/", postURLHandler.PostURL)
	e.POST("/api/shorten", postURLJSONHandler.PostURLJSON)
	if err := e.Start(cfg.SrvAddr); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
