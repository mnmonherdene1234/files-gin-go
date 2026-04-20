package main

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandleUploadEscapesDownloadURL(t *testing.T) {
	dir := t.TempDir()
	app := NewApp(Config{
		FilesDir:          dir,
		StaticFilesPath:   "/files",
		ServeStaticFiles:  true,
		MaxUploadMemoryMB: 32,
		MaxUploadSizeMB:   100,
	})

	req, err := newMultipartUploadRequest("report 2026 #1.txt", "hello", true)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}

	rr := httptest.NewRecorder()
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp UploadResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.Filename != "report 2026 #1.txt" {
		t.Fatalf("unexpected filename: %q", resp.Filename)
	}

	if resp.DownloadURL != "/files/report%202026%20%231.txt" {
		t.Fatalf("unexpected download URL: %q", resp.DownloadURL)
	}

	if _, err := os.Stat(filepath.Join(dir, "report 2026 #1.txt")); err != nil {
		t.Fatalf("expected uploaded file to exist: %v", err)
	}
}

func TestCORSIncludesConfiguredAPIKeyHeader(t *testing.T) {
	app := NewApp(Config{
		APIKeyEnabled:    true,
		APIKeyHeader:     "X-Custom-Key",
		APIKey:           "secret",
		FilesDir:         t.TempDir(),
		ServeStaticFiles: false,
		StaticFilesPath:  "/files",
	})

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	rr := httptest.NewRecorder()

	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rr.Code)
	}

	allowHeaders := rr.Header().Get("Access-Control-Allow-Headers")
	if !strings.Contains(allowHeaders, "X-Custom-Key") {
		t.Fatalf("expected allow headers to include custom API key header, got %q", allowHeaders)
	}
}

func newMultipartUploadRequest(filename, contents string, useOriginalFilename bool) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fileWriter, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := fileWriter.Write([]byte(contents)); err != nil {
		return nil, err
	}

	if useOriginalFilename {
		if err := writer.WriteField("useOriginalFilename", "true"); err != nil {
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}
