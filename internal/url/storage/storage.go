package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type LinkEntity struct {
	ID          string `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

var Links []LinkEntity

func SaveLinkFile(filename string, newLink LinkEntity) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0777)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("file open: %s", err.Error())
	}

	fileBuffer, err := io.ReadAll(file)
	if len(fileBuffer) != 0 {
		if err = json.Unmarshal(fileBuffer, &Links); err != nil {
			return fmt.Errorf("unmarshaling: %s", err.Error())
		}
	}

	Links = append(Links, newLink)

	bytes, err := json.MarshalIndent(Links, "", "\t")
	if err != nil {
		return fmt.Errorf("marshaling: %s", err.Error())
	}

	err = os.WriteFile(filename, bytes, 0644)
	if err != nil {
		return fmt.Errorf("write file: %s", err.Error())
	}

	return nil
}

func ShortLinkById(filename string, id string) (string, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0777)
	defer file.Close()
	if err != nil {
		return "", fmt.Errorf("open file, %s", err.Error())
	}

	fileBuff, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("read file, %s", err.Error())
	}

	var links []LinkEntity
	err = json.Unmarshal(fileBuff, &links)
	if err != nil {
		return "", fmt.Errorf("unmarshaling: %s", err.Error())
	}

	for _, v := range links {
		if v.ID == id {
			return v.ShortURL, nil
		}
	}

	return "", errors.New("link not found")
}
