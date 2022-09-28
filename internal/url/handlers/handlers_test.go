package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type errorMessage struct {
	Message string `json:"message"`
}

func TestLink(t *testing.T) {
	testsPositive := []struct {
		name          string
		originalURL   string
		shortURL      string
		statusCode    int
		statusCodeGet int
		postRequest   string
		getRequest    string
	}{
		{
			name:          "LinkHandler google.com",
			originalURL:   "https://google.com",
			shortURL:      "http://localhost:8080/aHR0cHM6Ly9nb29nbGUuY29t",
			statusCode:    http.StatusCreated,
			statusCodeGet: http.StatusTemporaryRedirect,
			postRequest:   "http://localhost:8080",
			getRequest:    "http://localhost:8080/aHR0cHM6Ly9nb29nbGUuY29t",
		},
		{
			name:          "LinkHandler amazon.com",
			originalURL:   "https://amazon.com",
			shortURL:      "http://localhost:8080/aHR0cHM6Ly9hbWF6b24uY29t",
			statusCode:    http.StatusCreated,
			statusCodeGet: http.StatusTemporaryRedirect,
			postRequest:   "http://localhost:8080",
			getRequest:    "http://localhost:8080/aHR0cHM6Ly9hbWF6b24uY29t",
		},
	}

	testsNegative := []struct {
		name          string
		originalURL   string
		shortURL      string
		statusCode    int
		statusCodeGet int
		postRequest   string
		getRequest    string
		error         string
	}{
		{
			name:          "LinkHandler google.com",
			originalURL:   "https://google.com",
			shortURL:      "http://localhost:8080/aHR0cHM6Ly9nb29nbGUuY29t",
			statusCode:    http.StatusCreated,
			statusCodeGet: http.StatusNotFound,
			postRequest:   "http://localhost:8080",
			getRequest:    "http://localhost:8080/aHR0cHM6Ly9nb29",
			error:         "link not found",
		},
		{
			name:          "LinkHandler google.com",
			originalURL:   "https://google.com",
			shortURL:      "http://localhost:8080/aHR0cHM6Ly9nb29nbGUuY29t",
			statusCode:    http.StatusCreated,
			statusCodeGet: http.StatusBadRequest,
			postRequest:   "http://localhost:8080",
			getRequest:    "http://localhost:8080",
			error:         "id not found on postRequest",
		},
	}
	for _, test := range testsPositive {
		t.Run(test.name, func(t *testing.T) {
			bodyReader := strings.NewReader(test.originalURL)
			request := httptest.NewRequest(http.MethodPost, test.postRequest, bodyReader)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(LinkHandler)
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, test.statusCode, result.StatusCode)

			shortLinkResult, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.shortURL, string(shortLinkResult))

			requestGet := httptest.NewRequest(http.MethodGet, test.getRequest, nil)
			w = httptest.NewRecorder()
			h.ServeHTTP(w, requestGet)
			result = w.Result()

			assert.Equal(t, test.statusCodeGet, result.StatusCode)
			assert.Equal(t, test.originalURL, result.Header.Get("Location"))
		})
	}

	for _, test := range testsNegative {
		t.Run(test.name, func(t *testing.T) {
			requestGet := httptest.NewRequest(http.MethodGet, test.getRequest, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(LinkHandler)
			h.ServeHTTP(w, requestGet)
			result := w.Result()

			assert.Equal(t, test.statusCodeGet, result.StatusCode)

			responseBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			err = result.Body.Close()
			require.NoError(t, err)

			var errorMsg errorMessage

			err = json.Unmarshal(responseBody, &errorMsg)
			assert.Equal(t, test.error, errorMsg.Message)
		})
	}

	t.Run("Body is nil", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "https://yahoo.com", nil)
		w := httptest.NewRecorder()
		h := http.HandlerFunc(LinkHandler)
		h.ServeHTTP(w, request)
		result := w.Result()
		defer result.Body.Close()
		fmt.Println("STATUS", result.StatusCode)
		assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	})
}
