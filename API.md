# API Reference

## Authentication

When `API_KEY_ENABLED=true`, protected endpoints require:

```
X-API-Key: <your-api-key>
```

Public endpoints (`GET /`, `GET /health`) do not require the key.

---

## Endpoints

### GET /

Returns service info and available endpoints.

**Response** `200`

```json
{
  "name": "files-gin-go",
  "endpoints": ["GET /health", "POST /upload", "DELETE /delete", "GET /list", "GET /size"],
  "static_files_path": "/files"
}
```

---

### GET /health

Health check.

**Response** `200`

```json
{ "status": "ok" }
```

---

### POST /upload

Upload a file. Requires `Content-Type: multipart/form-data`.

**Form field:** `file` — the file to upload

**Form field:** `useOriginalFilename` _(optional)_ — set to `true` to keep the original filename. Default generates a unique name with a timestamp.

**Response** `200`

```json
{
  "message": "File uploaded successfully",
  "filename": "document_2026-04-20T10-30-00.000000000.pdf",
  "downloadUrl": "/files/document_2026-04-20T10-30-00.000000000.pdf"
}
```

`downloadUrl` is only included when static file serving is enabled.

**Errors**

- `400` — No file received
- `400` — Upload size exceeds limit
- `401` — Missing or invalid API key
- `409` — File already exists (when `useOriginalFilename=true`)

---

### DELETE /delete

Delete a file.

**Request body**

```json
{ "filename": "document.pdf" }
```

**Response** `200`

```json
{ "message": "File deleted successfully", "filename": "document.pdf" }
```

**Errors**

- `400` — Filename is required
- `401` — Missing or invalid API key
- `404` — File not found

---

### GET /list

List all stored files.

**Response** `200`

```json
{
  "files": [
    { "name": "document.pdf", "size": 1024 },
    { "name": "image.png", "size": 2048 }
  ]
}
```

**Errors**

- `401` — Missing or invalid API key

---

### GET /size

Get total size of all stored files.

**Response** `200`

```json
{ "size": 3072 }
```

**Errors**

- `401` — Missing or invalid API key

---

### GET /files/{filename}

Serve stored files statically. Only available when `IS_SERVE_STATIC_FILES=true`.

**Response** — the file with `Content-Type` determined by the server.

**Errors**

- `401` — Missing or invalid API key
- `404` — File not found

---

## Error Responses

All errors return a JSON body:

```json
{ "error": "descriptive message" }
```

| Status | Meaning |
|---|---|
| `400` | Bad request — invalid filename, missing file, etc. |
| `401` | Unauthorized — API key missing or invalid |
| `404` | Not found — file does not exist |
| `409` | Conflict — file already exists |
| `413` | Request entity too large — upload exceeds `MAX_UPLOAD_SIZE_MB` |
| `500` | Internal server error |

---

## CORS

All responses include:

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Accept, X-API-Key
```

Preflight `OPTIONS` requests are handled and return `204 No Content`.
