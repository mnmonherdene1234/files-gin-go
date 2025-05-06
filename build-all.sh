#!/bin/bash

APP_NAME=filesgingo

echo "🔨 Building $APP_NAME for Windows..."
GOOS=windows GOARCH=amd64 go build -o build/${APP_NAME}.exe main.go

echo "🔨 Building $APP_NAME for Linux..."
GOOS=linux GOARCH=amd64 go build -o build/${APP_NAME}-linux main.go

echo "🔨 Building $APP_NAME for macOS..."
GOOS=darwin GOARCH=amd64 go build -o build/${APP_NAME}-mac main.go

echo "✅ Done. Binaries are in the 'build/' directory."
