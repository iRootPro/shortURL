package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/irootpro/shorturl/internal/url/service"
	"github.com/irootpro/shorturl/internal/url/storage"
	"github.com/irootpro/shorturl/internal/url/usecases"
)

type RequestPOST struct {
	URL string `json:"url"`
}

type ResponsePOST struct {
	Result string `json:"result"`
}

type GetURLHandler struct {
	cfg *service.ConfigVars
}

type PostURLHandler struct {
	cfg *service.ConfigVars
}

type PostURLJSONHandler struct {
	cfg *service.ConfigVars
}

func NewGetURLHandler(cfg *service.ConfigVars) *GetURLHandler {
	return &GetURLHandler{
		cfg: cfg,
	}
}

func NewPostURLHandler(cfg *service.ConfigVars) *PostURLHandler {
	return &PostURLHandler{
		cfg: cfg,
	}
}

func NewPostURLJSONHandler(cfg *service.ConfigVars) *PostURLJSONHandler {
	return &PostURLJSONHandler{
		cfg: cfg,
	}
}

func (h *GetURLHandler) GetURL(c echo.Context) error {
	id := c.Param("hash")
	if id == "" {
		return c.String(http.StatusBadRequest, "id not found on postRequest")
	}

	storageFile := h.cfg.StoragePath
	if storageFile == "" {
		shortURL, err := service.ShortURLByID(storage.Links, id)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}

		c.Response().Header().Set("Location", shortURL)
		return c.String(http.StatusTemporaryRedirect, "")
	}

	shortURL, err := storage.ShortLinkByID(storageFile, id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	c.Response().Header().Set("Location", shortURL)
	return c.String(http.StatusTemporaryRedirect, "")

}

func (h *PostURLHandler) PostURL(c echo.Context) error {
	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil || string(body) == "" {
		return c.String(http.StatusBadRequest, "error read body from postRequest")
	}

	id := usecases.GenerateShortLink(body)
	link := storage.LinkEntity{
		ID:          id,
		OriginalURL: string(body),
		ShortURL:    fmt.Sprintf("%s/%s", h.cfg.BaseURL, id),
	}

	storageFile := h.cfg.StoragePath

	if storageFile == "" {
		storage.Links = append(storage.Links, link)
	} else {
		if err := storage.SaveLinkFile(storageFile, link); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	return c.String(http.StatusCreated, link.ShortURL)
}

func (h *PostURLJSONHandler) PostURLJSON(c echo.Context) error {
	var request RequestPOST

	if err := c.Bind(&request); err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	id := usecases.GenerateShortLink([]byte(request.URL))
	link := storage.LinkEntity{
		ID:          id,
		OriginalURL: request.URL,
		ShortURL:    fmt.Sprintf("%s/%s", h.cfg.BaseURL, id),
	}

	storageFile := h.cfg.StoragePath

	if storageFile == "" {
		storage.Links = append(storage.Links, link)
	} else {
		if err := storage.SaveLinkFile(storageFile, link); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	response := &ResponsePOST{
		Result: link.ShortURL,
	}

	_, err := json.Marshal(response)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusCreated, response)
}
