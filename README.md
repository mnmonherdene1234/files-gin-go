# FilePocket

Simple file storage API written in Go.

## Run

```bash
cp .env.example .env
go run .
```

FilePocket listens on the port from `.env` or the default `9935`.

## Run With Docker Compose

The easiest way to run FilePocket is with the prebuilt Docker image:
[`mnmonherdene/gofilepocket`](https://hub.docker.com/r/mnmonherdene/gofilepocket).

In any deployment folder, create a `docker-compose.yml` file:

```yaml
services:
  gofilepocket:
    image: mnmonherdene/gofilepocket:latest
    container_name: gofilepocket
    ports:
      - "9935:9935"
    environment:
      SERVER_PORT: "9935"
      FILES_DIR: /app/files
      STATIC_FILES_SERVE_PATH: /files
      IS_SERVE_STATIC_FILES: "true"
      API_KEY_ENABLED: "false"
      API_KEY_HEADER: X-API-Key
      API_KEY: ""
      MAX_UPLOAD_MEMORY_MB: "32"
    volumes:
      - gofilepocket-files:/app/files
    restart: unless-stopped

volumes:
  gofilepocket-files:
```

Start the service:

```bash
docker compose up -d
```

The API will be available at:

```text
http://localhost:9935
```

Uploaded files are stored in the Docker volume named `gofilepocket-files`.
This keeps files available after the container is stopped, restarted, or
recreated.

To stop the service:

```bash
docker compose down
```

To update to the newest image:

```bash
docker compose pull
docker compose up -d
```

### Enable API Key Protection

Set `API_KEY_ENABLED` to `true` and provide an `API_KEY`:

```yaml
environment:
  SERVER_PORT: "9935"
  FILES_DIR: /app/files
  API_KEY_ENABLED: "true"
  API_KEY_HEADER: X-API-Key
  API_KEY: change-this-secret
```

Requests to protected endpoints must include the configured header:

```bash
curl -H "X-API-Key: change-this-secret" http://localhost:9935/list
```

## Build Locally With Docker

If you are developing the app and want to build the image from this repository,
use the included `docker-compose.yml`:

```bash
docker compose up --build
```

The container listens on the port from `SERVER_PORT`, defaulting to `9935`.
Files are persisted in the `gofilepocket-files` volume at `/app/files`.

To change runtime settings, create a `.env` file in the project root before
starting Compose. The same `SERVER_PORT` value is used for both the host port
and the container port, so mapping stays consistent.
Compose uses `/app/files` inside the container for file storage. If you need a
different storage path, edit `docker-compose.yml` and the volume mount together.

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
| `MAX_UPLOAD_SIZE_MB`      | `0`         | Optional upload size limit in MB. `0` means unlimited |

When `API_KEY_ENABLED=true`, `API_KEY` must be set.
Uploads are unlimited by default. Set `MAX_UPLOAD_SIZE_MB` to a positive number
only when you want to enforce a total upload size limit.

## API

See [API.md](API.md) for the endpoint reference.
