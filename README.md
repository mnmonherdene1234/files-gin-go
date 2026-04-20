# FilePocket

Simple file storage API written in Go.

## Run

```bash
cp .env.example .env
go run .
```

FilePocket listens on the port from `.env` or the default `9935`.

## Run With Docker

Build and run the service with Docker Compose:

```bash
docker compose up --build
```

The container listens on the port from `SERVER_PORT`, defaulting to `9935`.
Files are persisted in the `gofilepocket-files` volume at `/app/files`.

To change runtime settings, create a `.env` file in the project root before
starting Compose. The same `SERVER_PORT` value is used for both the host port
and the container port, so mapping stays consistent.
If you change `FILES_DIR`, update the volume mount in `docker-compose.yml` to
match.

If you prefer plain Docker, build and run the image directly:

```bash
docker build -t gofilepocket .
docker run --rm -p 9935:9935 -v gofilepocket-files:/app/files gofilepocket
```

## Environment Variables

| Variable                  | Default     | Description                                |
| ------------------------- | ----------- | ------------------------------------------ |
| `SERVER_PORT`             | `9935`      | HTTP server port                           |
| `FILES_DIR`               | `./files`   | Local directory used to store files        |
| `STATIC_FILES_SERVE_PATH` | `/files`    | Public path for static file serving        |
| `IS_SERVE_STATIC_FILES`   | `true`      | Enable or disable static file serving      |
| `API_KEY_ENABLED`         | `false`     | Require an API key for protected endpoints |
| `API_KEY_HEADER`          | `X-API-Key` | Request header name for the API key        |
| `API_KEY`                 | empty       | API key value used when auth is enabled    |
| `MAX_UPLOAD_MEMORY_MB`    | `32`        | Memory threshold for multipart parsing     |
| `MAX_UPLOAD_SIZE_MB`      | `100`       | Total upload size limit in MB              |

When `API_KEY_ENABLED=true`, `API_KEY` must be set.

## API

See [API.md](API.md) for the endpoint reference.
