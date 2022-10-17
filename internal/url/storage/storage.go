package storage

type LinkEntity struct {
	ID          string `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

var Links []LinkEntity
