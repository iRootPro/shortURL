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

type StorageFile struct {
	filename string
}

type StorageMemory struct {
	links []LinkEntity
}

func NewStorageFile(filename string) *StorageFile {
	return &StorageFile{
		filename: filename,
	}
}

func NewStorageMemory() *StorageMemory {
	return &StorageMemory{
		links: []LinkEntity{},
	}
}

func (s *StorageFile) Put(newLink LinkEntity) error {
	file, err := os.OpenFile(s.filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("file open: %s", err.Error())
	}
	defer file.Close()

	fileBuffer, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("read from file: %s", err)
	}

	var links []LinkEntity

	if len(fileBuffer) != 0 {
		if err = json.Unmarshal(fileBuffer, &links); err != nil {
			return fmt.Errorf("unmarshaling: %s", err.Error())
		}
	}

	links = append(links, newLink)

	bytes, err := json.MarshalIndent(links, "", "\t")
	if err != nil {
		return fmt.Errorf("marshaling: %s", err.Error())
	}

	err = os.WriteFile(s.filename, bytes, 0644)
	if err != nil {
		return fmt.Errorf("write file: %s", err.Error())
	}

	return nil
}

func (s *StorageFile) Get(id string) (string, error) {
	file, err := os.OpenFile(s.filename, os.O_RDONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("open file, %s", err.Error())
	}
	defer file.Close()

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
			return v.OriginalURL, nil
		}
	}

	return "", errors.New("link not found")
}

func (s *StorageMemory) Put(link LinkEntity) error {
	s.links = append(s.links, link)
	return nil
}

func (s *StorageMemory) Get(id string) (string, error) {
	for _, v := range s.links {
		if v.ID == id {
			return v.OriginalURL, nil
		}
	}

	return "", fmt.Errorf("link not found")
}
