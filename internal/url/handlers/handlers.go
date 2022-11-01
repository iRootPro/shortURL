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

type ServerHandler struct {
	cfg     *service.ConfigVars
	storage Storage
}

func NewServerHandler(cfg *service.ConfigVars, storage Storage) *ServerHandler {
	return &ServerHandler{
		cfg:     cfg,
		storage: storage,
	}
}

type Storage interface {
	Put(newLink storage.LinkEntity) error
	Get(id string) (string, error)
	GetAll() ([]storage.LinkEntity, error)
	Close() error
}

func (h *ServerHandler) GetURL(c echo.Context) error {
	id := c.Param("hash")
	if id == "" {
		return c.String(http.StatusBadRequest, "id not found on postRequest")
	}

	shortURL, err := h.storage.Get(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	c.Response().Header().Set("Location", shortURL)
	return c.String(http.StatusTemporaryRedirect, "")
}

func (h *ServerHandler) GetURLs(c echo.Context) error {
	coockie, err := c.Cookie("token")
	if err != nil || !service.CheckCookie(coockie) {
		c.Response().Header().Set("Content-Type", "application/json")
		c.Response().WriteHeader(http.StatusNoContent)
		return nil
	}

	c.SetCookie(coockie)
	urls, err := h.storage.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, "error read from storage")
	}

	if len(urls) == 0 {
		c.Response().Header().Set("Content-Type", "application/json")
		c.Response().WriteHeader(http.StatusNoContent)
		return nil
	}

	bytes, err := json.Marshal(urls)
	if err != nil {
		return c.String(http.StatusInternalServerError, "error marhaling data")
	}

	return c.JSONBlob(http.StatusOK, bytes)
}

func (h *ServerHandler) PostURL(c echo.Context) error {
  cookie, err := c.Cookie("token")
	if err != nil || service.CheckCookie(cookie) {
		cookie = service.SetCookie()
	}

  c.SetCookie(cookie)

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

	if err := h.storage.Put(link); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusCreated, link.ShortURL)
}

func (h *ServerHandler) PostURLJSON(c echo.Context) error {
  cookie, err := c.Cookie("token")
	if err != nil || service.CheckCookie(cookie) {
		cookie = service.SetCookie()
	}

  c.SetCookie(cookie)

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

	if err := h.storage.Put(link); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	response := &ResponsePOST{
		Result: link.ShortURL,
	}

	_, err = json.Marshal(response)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusCreated, response)
}
