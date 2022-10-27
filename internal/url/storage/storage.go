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
  ID          string `json:"-"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

type StorageFile struct {
	file *os.File
  memory StorageMemory 
}

type StorageMemory struct {
	links []LinkEntity
}

func NewStorageFile(filename string) (*StorageFile, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("file open: %s", err.Error())
	}

  fileBuff, err := io.ReadAll(file)
  if err != nil {
    log.Fatal("read from file:", err)
  }

	var links []LinkEntity

	if len(fileBuff) != 0 {
		if err = json.Unmarshal(fileBuff, &links); err != nil {
			return nil, fmt.Errorf("unmarshaling: %s", err.Error())
		}
	}

  memory := NewStorageMemory()
	return &StorageFile{
		file: file,
    memory: *memory,
	}, nil
}

func NewStorageMemory() *StorageMemory {
	return &StorageMemory{
		links: []LinkEntity{},
	}
}

func (s *StorageFile) Put(newLink LinkEntity) error {
  s.memory.links= append(s.memory.links, newLink)
	return nil
}

func (s *StorageFile) Get(id string) (string, error) {
  fmt.Print(s.memory.links)
	for _, v := range s.memory.links {
		if v.ID == id {
			return v.OriginalURL, nil
		}
	}

	return "", errors.New("link not found")
}

func (s *StorageFile) GetAll() ([]LinkEntity, error) {
  return s.memory.links, nil
}

func (s *StorageFile) Close() error {
  fmt.Println("Save data to file")

  bytes, err := json.Marshal(s.memory.links)

  if err != nil {
    return fmt.Errorf("marshaling: %s", err.Error())
  }

  _, err = s.file.Write(bytes)
  if err != nil {
    return fmt.Errorf("write to file, when close file description: %s ", err.Error())
  }
	err = s.file.Close()
	if err != nil {
    return fmt.Errorf("close file descriptor: %s", err.Error())
	}

  return nil
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

func (s *StorageMemory) GetAll() ([]LinkEntity, error) {
  return s.links, nil
}

func (s *StorageMemory) Close() error {
	s.links = []LinkEntity{}
  return nil
}
