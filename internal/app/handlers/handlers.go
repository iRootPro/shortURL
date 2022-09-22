package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type LinkEntity struct {
	ID          string `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

var links []LinkEntity

func Link(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := strings.TrimPrefix(r.URL.Path, "/")
		if id == "" {
			errorResponse(w, "id not found on request", http.StatusBadRequest)
			return
		}

		for _, v := range links {
			if v.ID == id {
				w.Header().Add("Location", v.OriginalURL)
				w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}
		}
		errorResponse(w, "link not fount", http.StatusNotFound)
	case http.MethodPost:
		fmt.Println("post")
		if r.Header.Get("Content-type") != "application/json" {
			errorResponse(w, "Content-type is not application/json", http.StatusUnsupportedMediaType)
			return
		}
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			errorResponse(w, "Error read body from request", http.StatusBadRequest)
			return
		}
		fmt.Println(string(body))
		var link LinkEntity
		if err = json.Unmarshal(body, &link); err != nil {
			errorResponse(w, "Error unmarshaling error", http.StatusBadRequest)
			return
		}
		links = append(links, link)
		w.WriteHeader(http.StatusCreated)
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
