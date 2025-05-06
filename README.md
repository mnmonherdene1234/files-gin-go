# Files GIN Go

## Overview

This project is a file management API built using the GIN framework in Go. It provides endpoints for uploading, listing, deleting files, and calculating folder sizes.

## Getting Started

### Setup

1. Copy the `.env.example` file to `.env`:
   ```bash
   cp .env.example .env
   ```
2. Update the `.env` file with your own configuration as needed.

### Running the Application

To start the application, use the provided `run.sh` script to build and run the source code:

```bash
./run.sh
```

Alternatively, you can use the pre-built release binary if available.

## Swagger (Development Only)

Swagger is a tool that helps you document and test your APIs. This section is intended for development purposes only.

### Swagger UI Endpoint

You can access the Swagger UI at the following endpoint:

```
/swagger/index.html
```

### Installation

To install the Swagger CLI tool, run the following command:

```
go install github.com/swaggo/swag/cmd/swag@latest
```

### Initialize Swagger

To generate the Swagger documentation for your project, run:

```
swag init
```
