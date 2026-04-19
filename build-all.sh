#!/bin/bash

APP_NAME=filesgingo

mkdir -p build

GOOS=windows GOARCH=amd64 go build -o build/${APP_NAME}.exe .
GOOS=linux GOARCH=amd64 go build -o build/${APP_NAME}-linux .
GOOS=darwin GOARCH=amd64 go build -o build/${APP_NAME}-mac .
