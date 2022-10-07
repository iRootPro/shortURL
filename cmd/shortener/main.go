package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/irootpro/shorturl/internal/url/handlers"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/:hash", handlers.GetURL)
	e.POST("/", handlers.PostURL)
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
