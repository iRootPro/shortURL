package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	e.POST("/api/shorten/batch", serverHandler.PostURLsBatchJSON)
	e.GET("/api/user/urls", serverHandler.GetURLs)
	e.DELETE("/api/user/urls", serverHandler.RemoveURLs)
	e.GET("/ping", serverHandler.Ping)
	go func() {
		if err := e.Start(cfg.SrvAddr); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer storage.Close()
	defer fmt.Println("Server shutdown")
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
