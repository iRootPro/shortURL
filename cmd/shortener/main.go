package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/irootpro/shorturl/internal/url/handlers"
)

func main() {
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "localhost:8080"
	}
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/:hash", handlers.GetURL)
	e.POST("/", handlers.PostURL)
	e.POST("/api/shorten", handlers.PostURLJSON)
	if err := e.Start(serverAddress); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
