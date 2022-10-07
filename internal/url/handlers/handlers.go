package handlers

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/irootpro/shorturl/internal/url/storage"
	"github.com/irootpro/shorturl/internal/url/usecases"
)

const host = "http://localhost:8080"

func GetURL(c echo.Context) error {
	id := c.Param("hash")
	if id == "" {
		return c.String(http.StatusBadRequest, "id not found on postRequest")
	}

	for _, v := range storage.Links {
		if v.ID == id {
			c.Response().Header().Set("Location", v.OriginalURL)
			return c.String(http.StatusTemporaryRedirect, "")
		}
	}

	return c.String(http.StatusNotFound, "link not found")
}

func PostURL(c echo.Context) error {
	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil || string(body) == "" {
		return c.String(http.StatusBadRequest, "error read body from postRequest")
	}

	id := usecases.GenerateShortLink(body)
	link := storage.LinkEntity{
		ID:          id,
		OriginalURL: string(body),
		ShortURL:    host + "/" + id,
	}
	storage.Links = append(storage.Links, link)
	return c.String(http.StatusCreated, link.ShortURL)
}
