package utils

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func Createdb(s string) {
	f, err := isValidDirectoryName(s)
	if err != nil {
		fmt.Printf("Failed to initialize data files: %v\n", err)
		os.Exit(1)
	}
	if f {
		if err := InitializeDB(s); err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			os.Exit(1)
		}
	}
}

// Проверка валидности имени директории
func isValidDirectoryName(name string) (bool, error) {
	reservedNames := map[string]bool{
		"cmd":       true,
		"internal":  true,
		"models":    true,
		"go.mod":    true,
		"README.md": true,
		// filepath.Join("cmd", "hot-coffee"): true,
	}

	if reservedNames[name] {
		return false, errors.New("the chosen directory name is reserved and cannot be used")
	}

	if strings.TrimSpace(name) == "" {
		return false, errors.New("directory name cannot be empty or contain only whitespace")
	}

	if net.ParseIP(name) != nil {
		return false, errors.New("directory name must not be an IP address")
	}

	if filepath.IsAbs(name) || filepath.Clean(name) != name {
		return false, errors.New("invalid directory path, must be relative to project root")
	}

	validDirPattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	if !validDirPattern.MatchString(name) {
		return false, errors.New("directory name can only contain letters, numbers, dots, underscores, and hyphens")
	}

	return true, nil
}

// Создание папки в корне проекта
func InitializeDB(dataDir string) error {
	dataPath := filepath.Join("..", dataDir)

	// If directory exists — remove it
	if _, err := os.Stat(dataPath); err == nil {
		if err := os.RemoveAll(dataPath); err != nil {
			return fmt.Errorf("failed to remove old directory: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check directory status: %w", err)
	}

	// Create new data directory
	if err := os.MkdirAll(dataPath, 0o766); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create valid empty JSON files with []
	for _, fileName := range []string{"orders.json", "menu_items.json", "inventory.json"} {
		filePath := filepath.Join(dataPath, fileName)
		if err := os.WriteFile(filePath, []byte("[]"), 0o644); err != nil {
			return fmt.Errorf("failed to initialize file %s: %w", fileName, err)
		}
	}

	slog.Info("Data directory and files initialized", "path", dataPath)
	return nil
}
