package main

import (
	"errors"
	"os"
	"path/filepath"
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
