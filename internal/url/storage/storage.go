package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	apiError "github.com/irootpro/shorturl/internal/error"
	"github.com/irootpro/shorturl/internal/url/usecases"
	"io"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

type LinkEntity struct {
	ID          string `json:"-"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

type LinkBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type LinkBatchResult struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
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
		"CREATE TABLE IF NOT EXISTS links (hash_url TEXT NOT NULL, original_url TEXT NOT NULL UNIQUE, short_url TEXT NOT NULL, correlation_id TEXT, is_deleted VARCHAR(10) DEFAULT '')",
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

func (s *StorageMemory) Batch(ctx context.Context, links []LinkBatch, baseURL string) ([]LinkBatchResult, error) {
	return []LinkBatchResult{}, nil
}

func (s *StorageFile) Batch(ctx context.Context, links []LinkBatch, baseURL string) ([]LinkBatchResult, error) {
	return []LinkBatchResult{}, nil
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

func (s *StorageFile) RemoveURLs(urls []string) error {
	if len(urls) == 0 {
		return nil
	}

	var afterRemovedLinks []LinkEntity
	for _, link := range s.memory.links {
		for _, url := range urls {
			if link.ShortURL != url {
				afterRemovedLinks = append(afterRemovedLinks, link)
			}
		}
	}

	s.memory.links = afterRemovedLinks

	return nil
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

func (s *StorageMemory) RemoveURLs(urls []string) error {
	if len(urls) == 0 {
		return nil
	}

	var afterRemovedLinks []LinkEntity
	for _, link := range s.links {
		for _, url := range urls {
			if link.ShortURL != url {
				afterRemovedLinks = append(afterRemovedLinks, link)
			}
		}
	}

	s.links = afterRemovedLinks
	return nil
}

func (s *StorageMemory) Ping() error {
	return nil
}

func (s *StorageDB) Put(link LinkEntity) error {
	row, err := s.db.Query("INSERT INTO links VALUES ($1, $2, $3) ", link.ID, link.OriginalURL, link.ShortURL)
	if err != nil {
		return &apiError.NotUniqueRecordError{
			URL: link.OriginalURL,
		}
	}

	if err = row.Err(); err != nil {
		return fmt.Errorf("row err: %s", err.Error())
	}

	defer row.Close()
	return nil
}

func (s *StorageDB) Get(id string) (string, error) {
	row := s.db.QueryRow("SELECT original_url, is_deleted from links WHERE hash_url=$1", id)
	var originalURL string
	var isDeleted string
	if err := row.Scan(&originalURL, &isDeleted); err != nil {
		return "", fmt.Errorf("get url by id: %s", err.Error())
	}

	if isDeleted == "deleted" {
		return "deleted", nil
	}

	return originalURL, nil
}

func (s *StorageDB) Batch(ctx context.Context, links []LinkBatch, baseURL string) ([]LinkBatchResult, error) {
	result := make([]LinkBatchResult, 0)
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("start transaction, %s", err.Error())
	}

	defer tx.Rollback()

	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO links(hash_url, short_url, original_url, correlation_id) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return nil, fmt.Errorf("prepare statement, %s", err.Error())
	}

	defer stmt.Close()

	for _, v := range links {
		short := usecases.GenerateShortLink([]byte(v.OriginalURL))
		if _, err := stmt.ExecContext(ctx, short, fmt.Sprintf("%s/%s", baseURL, short), v.OriginalURL, v.CorrelationID); err != nil {
			return nil, fmt.Errorf("statement exec, %s", err.Error())
		}

		result = append(result, LinkBatchResult{
			CorrelationID: v.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", baseURL, short),
		})
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit, %s", err.Error())
	}

	return result, nil
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

func (s *StorageDB) RemoveURLs(urls []string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("start transaction, %s", err.Error())
	}

	if len(urls) == 0 {
		return errors.New("list of URLs is empty")
	}

	ctx := context.Background()

	stmt, err := tx.PrepareContext(ctx, "UPDATE links SET is_deleted='deleted' WHERE hash_url=$1")
	if err != nil {
		return fmt.Errorf("prepare statement, %s", err.Error())
	}
	defer stmt.Close()

	fanOutsChan := fanOut(urls, len(urls))
	wg := &sync.WaitGroup{}
	errCh := make(chan error)

	urlsChan := make(chan string, len(urls))

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() error {
			err = executer(ctx, stmt, tx, urlsChan)
			if err != nil {
				wg.Done()
				return err
			}
			wg.Done()
			return nil
		}()
	}

	for _, item := range fanOutsChan {
		urlsChan <- <-item
	}
	close(urlsChan)

	go func() {
		wg.Wait()
		close(errCh)
	}()

	if err = <-errCh; err != nil {
		return fmt.Errorf("error from channel, %s", err.Error())
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit db, %s", err.Error())
	}

	return nil
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

func fanOut(input []string, n int) []chan string {
	chs := make([]chan string, 0, n)
	for i, v := range input {
		ch := make(chan string, 1)
		ch <- v
		chs = append(chs, ch)
		close(chs[i])
	}
	return chs
}

func executer(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, inputChan <-chan string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for id := range inputChan {
		if _, err := stmt.ExecContext(ctx, id); err != nil {
			if err = tx.Rollback(); err != nil {
				return fmt.Errorf("rollback, %s", err.Error())
			}
			return fmt.Errorf("exec, %s", err.Error())
		}
	}
	return nil
}
