package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDeleteRejectsPathTraversal(t *testing.T) {
	store := NewFileStore(t.TempDir())

	err := store.Delete("../outside.txt")
	if !errors.Is(err, ErrInvalidFilename) {
		t.Fatalf("expected ErrInvalidFilename, got %v", err)
	}
}

func TestListAndFolderSize(t *testing.T) {
	dir := t.TempDir()
	store := NewFileStore(dir)

	content := []byte("hello")
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), content, 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	files, err := store.List()
	if err != nil {
		t.Fatalf("list files: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	if files[0].Name != "a.txt" {
		t.Fatalf("unexpected file name: %s", files[0].Name)
	}
	if files[0].Size != int64(len(content)) {
		t.Fatalf("unexpected file size: %d", files[0].Size)
	}

	size, err := store.FolderSize()
	if err != nil {
		t.Fatalf("folder size: %v", err)
	}
	if size != int64(len(content)) {
		t.Fatalf("unexpected folder size: %d", size)
	}
}

func TestExists(t *testing.T) {
	dir := t.TempDir()
	store := NewFileStore(dir)

	exists, err := store.Exists("nonexistent.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Fatal("expected nonexistent file to not exist")
	}

	if err := os.WriteFile(filepath.Join(dir, "test.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	exists, err = store.Exists("test.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Fatal("expected file to exist")
	}
}

func TestExistsRejectsPathTraversal(t *testing.T) {
	store := NewFileStore(t.TempDir())

	_, err := store.Exists("../outside.txt")
	if !errors.Is(err, ErrInvalidFilename) {
		t.Fatalf("expected ErrInvalidFilename, got %v", err)
	}
}

func TestSaveAtomicDuplicate(t *testing.T) {
	dir := t.TempDir()
	store := NewFileStore(dir)

	err := store.Save(strings.NewReader("hello"), "duplicate.txt")
	if err != nil {
		t.Fatalf("first save failed: %v", err)
	}

	err = store.Save(strings.NewReader("hello"), "duplicate.txt")
	if !errors.Is(err, ErrFileAlreadyExists) {
		t.Fatalf("expected ErrFileAlreadyExists, got %v", err)
	}
}

func TestSaveAtomicDuplicateRejectsPathTraversal(t *testing.T) {
	store := NewFileStore(t.TempDir())

	err := store.Save(strings.NewReader("test"), "../evil.txt")
	if !errors.Is(err, ErrInvalidFilename) {
		t.Fatalf("expected ErrInvalidFilename, got %v", err)
	}
}
