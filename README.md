# files-gin-go

Simple file management API written with Go standard library only.

## File Structure

```text
.
├── app.go
├── config.go
├── main.go
├── store.go
├── store_test.go
├── .env.example
├── go.mod
├── run.sh
└── build-all.sh
```

## Endpoints

- `POST /upload`
- `DELETE /delete`
- `GET /list`
- `GET /size`
- `GET /files/...` when `IS_SERVE_STATIC_FILES=true`

If `API_KEY_ENABLED=true`, upload, delete, list, size, and static file access require the header defined by `API_KEY_HEADER`.

## Run

```bash
cp .env.example .env
go run .
```

## Notes

- No external runtime dependency is used.
- `.env` loading is implemented with standard library code.
- CORS is handled manually.
