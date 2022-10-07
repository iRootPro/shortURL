package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLink(t *testing.T) {
	testsPositive := []struct {
		name          string
		originalURL   string
		shortURL      string
		statusCode    int
		statusCodeGet int
		postRequest   string
		hashShortURL  string
	}{
		{
			name:          "LinkHandler google.com",
			originalURL:   "https://google.com",
			shortURL:      "http://localhost:8080/aHR0cHM6Ly9nb29nbGUuY29t",
			statusCode:    http.StatusCreated,
			statusCodeGet: http.StatusTemporaryRedirect,
			postRequest:   "http://localhost:8080",
			hashShortURL:  "aHR0cHM6Ly9nb29nbGUuY29t",
		},
		{
			name:          "LinkHandler amazon.com",
			originalURL:   "https://amazon.com",
			shortURL:      "http://localhost:8080/aHR0cHM6Ly9hbWF6b24uY29t",
			statusCode:    http.StatusCreated,
			statusCodeGet: http.StatusTemporaryRedirect,
			postRequest:   "http://localhost:8080",
			hashShortURL:  "aHR0cHM6Ly9hbWF6b24uY29t",
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
			getRequest:    "aHR0cHM6Ly9nb29",
			error:         "link not found",
		},
		{
			name:          "LinkHandler google.com",
			originalURL:   "https://google.com",
			shortURL:      "http://localhost:8080/aHR0cHM6Ly9nb29nbGUuY29t",
			statusCode:    http.StatusCreated,
			statusCodeGet: http.StatusBadRequest,
			postRequest:   "http://localhost:8080",
			getRequest:    "",
			error:         "id not found on postRequest",
		},
	}
	for _, test := range testsPositive {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			bodyReader := strings.NewReader(test.originalURL)
			requestPost := httptest.NewRequest(http.MethodPost, test.postRequest, bodyReader)
			w := httptest.NewRecorder()
			c := e.NewContext(requestPost, w)
			PostURL(c)
			result := w.Result()

			assert.Equal(t, test.statusCode, result.StatusCode)

			shortLinkResult, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.shortURL, string(shortLinkResult))

			requestGet := httptest.NewRequest(http.MethodGet, "/", nil)
			w = httptest.NewRecorder()
			c = e.NewContext(requestGet, w)
			c.SetPath(":hash")
			c.SetParamNames("hash")
			c.SetParamValues(test.hashShortURL)
			GetURL(c)
			result = w.Result()
			defer result.Body.Close()
			assert.Equal(t, test.statusCodeGet, result.StatusCode)
			assert.Equal(t, test.originalURL, result.Header.Get("Location"))
		})
	}

	for _, test := range testsNegative {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			requestGet := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			c := e.NewContext(requestGet, w)
			c.SetPath(":hash")
			c.SetParamNames("hash")
			c.SetParamValues(test.getRequest)
			GetURL(c)
			result := w.Result()

			assert.Equal(t, test.statusCodeGet, result.StatusCode)

			responseBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			err = result.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, test.error, string(responseBody))
		})
	}
}
