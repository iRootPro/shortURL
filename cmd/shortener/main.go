package main

import (
	"net/http"

	"github.com/irootpro/shorturl/internal/url/handlers"
)

func main() {
	http.HandleFunc("/", handlers.Link)
	http.ListenAndServe(":8080", nil)
}
