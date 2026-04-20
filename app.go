package main

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type App struct {
	config Config
	store  *FileStore
}

type DeleteFileRequest struct {
	Filename string `json:"filename"`
}

type UploadResponse struct {
	Message     string `json:"message"`
	Filename    string `json:"filename"`
	DownloadURL string `json:"downloadUrl,omitempty"`
}

type DeleteResponse struct {
	Message  string `json:"message"`
	Filename string `json:"filename"`
}

type FilesListResponse struct {
	Files []FileMeta `json:"files"`
}

type SizeResponse struct {
	Size int64 `json:"size"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewApp(cfg Config) *App {
	return &App{
		config: cfg,
		store:  NewFileStore(cfg.FilesDir),
	}
}

func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", a.public(http.HandlerFunc(a.handleIndex)))
	mux.Handle("GET /health", a.public(http.HandlerFunc(a.handleHealth)))
	mux.Handle("POST /upload", a.protected(http.HandlerFunc(a.handleUpload)))
	mux.Handle("DELETE /delete", a.protected(http.HandlerFunc(a.handleDelete)))
	mux.Handle("GET /list", a.protected(http.HandlerFunc(a.handleList)))
	mux.Handle("GET /size", a.protected(http.HandlerFunc(a.handleSize)))

	if a.config.ServeStaticFiles {
		staticPrefix := a.config.StaticFilesPath + "/"
		fileServer := http.StripPrefix(staticPrefix, http.FileServer(http.Dir(a.config.FilesDir)))
		mux.Handle(staticPrefix, a.protected(fileServer))
		mux.Handle(a.config.StaticFilesPath, a.protected(http.RedirectHandler(staticPrefix, http.StatusMovedPermanently)))
	}

	return a.withCORS(a.withLogging(mux))
}

func (a *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"name": "FilePocket",
		"endpoints": []string{
			"GET /health",
			"POST /upload",
			"DELETE /delete",
			"GET /list",
			"GET /size",
		},
		"static_files_path": a.config.StaticFilesPath,
	})
}

func (a *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *App) handleUpload(w http.ResponseWriter, r *http.Request) {
	maxSize := a.config.MaxUploadSizeMB << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)
	if err := r.ParseMultipartForm(a.config.MaxUploadMemoryMB << 20); err != nil {
		if r.MultipartForm != nil {
			r.MultipartForm.RemoveAll()
		}
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			writeError(w, http.StatusRequestEntityTooLarge, "Upload size exceeds limit", err)
			return
		}
		writeError(w, http.StatusBadRequest, "Invalid multipart form", err)
		return
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "No file received", err)
		return
	}
	defer file.Close()

	filename := SafeBaseName(header.Filename)
	if !strings.EqualFold(r.FormValue("useOriginalFilename"), "true") {
		filename = UniqueFilename(filename)
	}

	if err := a.store.Save(file, filename); err != nil {
		writeStoreError(w, "Failed to save the file", err)
		return
	}

	var downloadURL string
	if a.config.ServeStaticFiles {
		downloadURL = a.config.StaticFilesPath + "/" + url.PathEscape(filename)
	}

	writeJSON(w, http.StatusOK, UploadResponse{
		Message:     "File uploaded successfully",
		Filename:    filename,
		DownloadURL: downloadURL,
	})
}

func (a *App) handleDelete(w http.ResponseWriter, r *http.Request) {
	var req DeleteFileRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	if strings.TrimSpace(req.Filename) == "" {
		writeError(w, http.StatusBadRequest, "Filename is required", errors.New("missing filename"))
		return
	}

	if err := a.store.Delete(req.Filename); err != nil {
		writeStoreError(w, "Failed to delete the file", err)
		return
	}

	writeJSON(w, http.StatusOK, DeleteResponse{
		Message:  "File deleted successfully",
		Filename: req.Filename,
	})
}

func (a *App) handleList(w http.ResponseWriter, r *http.Request) {
	files, err := a.store.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list files", err)
		return
	}

	writeJSON(w, http.StatusOK, FilesListResponse{Files: files})
}

func (a *App) handleSize(w http.ResponseWriter, r *http.Request) {
	size, err := a.store.FolderSize()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to calculate folder size", err)
		return
	}

	writeJSON(w, http.StatusOK, SizeResponse{Size: size})
}

func (a *App) public(next http.Handler) http.Handler {
	return next
}

func (a *App) protected(next http.Handler) http.Handler {
	if !a.config.APIKeyEnabled {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(a.config.APIKeyHeader)
		if apiKey == "" {
			writeError(w, http.StatusUnauthorized, "API key is required", errors.New("missing api key"))
			return
		}
		if subtle.ConstantTimeCompare([]byte(apiKey), []byte(a.config.APIKey)) != 1 {
			writeError(w, http.StatusUnauthorized, "Invalid API key", errors.New("invalid api key"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *App) withCORS(next http.Handler) http.Handler {
	allowHeaders := "Content-Type, Accept"
	if a.config.APIKeyHeader != "" {
		allowHeaders += ", " + a.config.APIKeyHeader
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *App) withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, recorder.status, time.Since(start))
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func writeStoreError(w http.ResponseWriter, fallbackMessage string, err error) {
	switch {
	case errors.Is(err, ErrInvalidFilename):
		writeError(w, http.StatusBadRequest, "Invalid filename", err)
	case errors.Is(err, ErrFileNotFound):
		writeError(w, http.StatusNotFound, "File not found", err)
	case errors.Is(err, ErrFileAlreadyExists):
		writeError(w, http.StatusConflict, "File already exists", err)
	default:
		writeError(w, http.StatusInternalServerError, fallbackMessage, err)
	}
}

func writeError(w http.ResponseWriter, status int, message string, err error) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
	writeJSON(w, status, ErrorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
