# files-gin-go

Simple file storage API using only Go standard library.

## Run

```bash
cp .env.example .env
go run .
```

Server starts on the port configured in `.env` (default `9935`).

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `SERVER_PORT` | `9935` | Port to listen on |
| `FILES_DIR` | `./files` | Directory for stored files |
| `STATIC_FILES_SERVE_PATH` | `/files` | URL path for static file serving |
| `IS_SERVE_STATIC_FILES` | `true` | Enable static file serving at `STATIC_FILES_SERVE_PATH` |
| `API_KEY_ENABLED` | `false` | Require API key for protected endpoints |
| `API_KEY_HEADER` | `X-API-Key` | HTTP header for the API key |
| `API_KEY` | _(empty)_ | The API key value |
| `MAX_UPLOAD_MEMORY_MB` | `32` | Memory buffer before spilling to temp files |
| `MAX_UPLOAD_SIZE_MB` | `100` | Max upload size in MB |

When `API_KEY_ENABLED=true`, the `API_KEY` must be set.

## API

See [API.md](API.md) for full endpoint documentation.
