package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/irootpro/shorturl/internal/url/storage"
	"github.com/irootpro/shorturl/internal/url/usecases"
)

const host = "http://localhost:8080"

func Link(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := strings.TrimPrefix(r.URL.Path, "/")
		if id == "" {
			errorResponse(w, "id not found on request", http.StatusBadRequest)
			return
		}

		for _, v := range storage.Links {
			if v.ID == id {
				w.Header().Add("Location", v.OriginalURL)
				w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}
		}
		errorResponse(w, "link not fount", http.StatusNotFound)
	case http.MethodPost:
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			errorResponse(w, "Error read body from request", http.StatusBadRequest)
			return
		}

		id := usecases.GenerateShortLink(body)
		link := storage.LinkEntity{
			ID:          id,
			OriginalURL: string(body),
			ShortURL:    host + "/" + id,
		}
		storage.Links = append(storage.Links, link)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(link.ShortURL))
	}
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(httpStatusCode)
	response := make(map[string]string)
	response["message"] = message
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}
