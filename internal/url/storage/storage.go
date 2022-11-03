package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type LinkEntity struct {
	ID          string `json:"-"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

type StorageFile struct {
	file   *os.File
	memory StorageMemory
}

type StorageMemory struct {
	links []LinkEntity
}

type StorageDB struct {
	db *sql.DB
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
		file:   file,
		memory: *memory,
	}, nil
}

func NewStorageMemory() *StorageMemory {
	return &StorageMemory{
		links: []LinkEntity{},
	}
}

func NewStorageDB(dsn string) *StorageDB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("connect to data base: %s", err.Error())
	}

	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS links (hash_url TEXT NOT NULL, original_url TEXT NOT NULL, short_url TEXT NOT NULL)",
	)
	if err != nil {
		log.Fatalf("create database: %s", err.Error())
	}

	return &StorageDB{
		db: db,
	}
}

func (s *StorageFile) Put(newLink LinkEntity) error {
	s.memory.links = append(s.memory.links, newLink)
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

func (s *StorageFile) Ping() error {
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

func (s *StorageMemory) Ping() error {
	return nil
}

func (s *StorageDB) Put(link LinkEntity) error {
	row, err := s.db.Query("INSERT INTO links VALUES ($1, $2, $3)", link.ID, link.OriginalURL, link.ShortURL)
	if err != nil {
		return fmt.Errorf("insert into database: %s", err.Error())
	}

	if err = row.Err(); err != nil {
		return fmt.Errorf("row err: %s", err.Error())
	}

	defer row.Close()
	return nil
}

func (s *StorageDB) Get(id string) (string, error) {
	row := s.db.QueryRow("SELECT original_url from links WHERE hash_url=$1", id)
	var originalURL string
	if err := row.Scan(&originalURL); err != nil {
		return "", fmt.Errorf("get url by id: %s", err.Error())
	}

	return originalURL, nil
}

func (s *StorageDB) GetAll() ([]LinkEntity, error) {
	rows, err := s.db.Query("SELECT * from links")
	if err != nil {
		return []LinkEntity{}, fmt.Errorf("get all urls: %s", err.Error())
	}

  if err = rows.Err(); err != nil {
    return []LinkEntity{}, fmt.Errorf("row scan: %s", err.Error())
  }
	defer rows.Close()

	var links []LinkEntity
	for rows.Next() {
		var link LinkEntity
		if err := rows.Scan(&link.ID, &link.OriginalURL, &link.ShortURL); err != nil {
			return links, fmt.Errorf("row scan: %s", err.Error())
		}
		links = append(links, link)
	}
	return links, nil
}

func (s *StorageDB) Close() error {
	return nil
}

func (s *StorageDB) Ping() error {
	if err := s.db.Ping(); err != nil {
		return fmt.Errorf("connect to database failed: %s", err.Error())
	}

	fmt.Println("Ping to database successful, connection is still alive")

	return nil
}
