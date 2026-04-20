package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrInvalidFilename  = errors.New("invalid filename")
	ErrFileNotFound     = errors.New("file not found")
	ErrFileAlreadyExists = errors.New("file already exists")
)

type FileMeta struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type FileStore struct {
	baseDir string
}

func NewFileStore(baseDir string) *FileStore {
	return &FileStore{baseDir: baseDir}
}

func SafeBaseName(name string) string {
	base := filepath.Base(strings.TrimSpace(name))
	if base == "" || base == "." {
		return "file"
	}
	return base
}

func UniqueFilename(originalName string) string {
	safeName := SafeBaseName(originalName)
	ext := filepath.Ext(safeName)
	baseName := strings.TrimSuffix(safeName, ext)
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000000")
	return fmt.Sprintf("%s_%s%s", baseName, timestamp, ext)
}

func (s *FileStore) Save(reader io.Reader, name string) error {
	filePath, err := s.resolveFilePath(name)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(s.baseDir, 0o755); err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return ErrFileAlreadyExists
		}
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

func (s *FileStore) Delete(name string) error {
	filePath, err := s.resolveFilePath(name)
	if err != nil {
		return err
	}

	if err := os.Remove(filePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrFileNotFound
		}
		return err
	}

	return nil
}

func (s *FileStore) Exists(name string) (bool, error) {
	filePath, err := s.resolveFilePath(name)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *FileStore) List() ([]FileMeta, error) {
	if err := os.MkdirAll(s.baseDir, 0o755); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, err
	}

	files := make([]FileMeta, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		files = append(files, FileMeta{
			Name: entry.Name(),
			Size: info.Size(),
		})
	}

	return files, nil
}

func (s *FileStore) FolderSize() (int64, error) {
	if err := os.MkdirAll(s.baseDir, 0o755); err != nil {
		return 0, err
	}

	var total int64
	err := filepath.WalkDir(s.baseDir, func(_ string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}
		total += info.Size()
		return nil
	})

	return total, err
}

func (s *FileStore) resolveFilePath(name string) (string, error) {
	normalizedName, err := normalizeStoredFilename(name)
	if err != nil {
		return "", err
	}

	baseDir, err := filepath.Abs(s.baseDir)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(baseDir, normalizedName)
	relativePath, err := filepath.Rel(baseDir, filePath)
	if err != nil {
		return "", err
	}
	if relativePath == ".." || strings.HasPrefix(relativePath, ".."+string(filepath.Separator)) {
		return "", ErrInvalidFilename
	}

	return filePath, nil
}

func normalizeStoredFilename(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" || trimmed == "." {
		return "", ErrInvalidFilename
	}
	if trimmed != filepath.Base(trimmed) {
		return "", ErrInvalidFilename
	}
	if strings.Contains(trimmed, "/") || strings.Contains(trimmed, "\\") {
		return "", ErrInvalidFilename
	}
	return trimmed, nil
}
