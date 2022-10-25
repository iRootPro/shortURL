package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type LinkEntity struct {
	ID          string `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

type StorageFile struct {
	file *os.File
}

type StorageMemory struct {
	links []LinkEntity
}

func NewStorageFile(filename string) *StorageFile {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("file open: %s", err.Error())
	}
	return &StorageFile{
		file: file,
	}
}

func NewStorageMemory() *StorageMemory {
	return &StorageMemory{
		links: []LinkEntity{},
	}
}

func (s *StorageFile) Put(newLink LinkEntity) error {
	fileBuffer, err := io.ReadAll(s.file)
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

	_, err = s.file.Write(bytes)
	if err != nil {
		return fmt.Errorf("write file: %s", err.Error())
	}

	return nil
}

func (s *StorageFile) Get(id string) (string, error) {

	fileBuff, err := io.ReadAll(s.file)
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

func (s *StorageFile) Close() {
	err := s.file.Close()
	if err != nil {
		return
	}
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

func (s *StorageMemory) Close() {
	s.links = []LinkEntity{}
}
