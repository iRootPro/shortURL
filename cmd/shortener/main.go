package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/irootpro/shorturl/internal/url/handlers"
)

func main() {
	http.HandleFunc("/", handlers.LinkHandler)
	fmt.Printf("Starting server at port 8080\n")
	log.Fatalln(http.ListenAndServe(":8080", nil))

}
