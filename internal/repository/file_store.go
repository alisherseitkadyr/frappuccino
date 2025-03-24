package repository

import (
	"encoding/json"
	"os"
	"sync"
)

type FileStore struct {
	mu       sync.Mutex
	filePath string
}

func NewFileStore(filePath string) *FileStore {
	return &FileStore{filePath: filePath}
}

func (fs *FileStore) Read(data interface{}) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	file, err := os.ReadFile(fs.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, data)
}

func (fs *FileStore) Write(data interface{}) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fs.filePath, file, 0644)
}