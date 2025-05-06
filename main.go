package main

import (
	_ "github.com/mnmonherdene1234/files-gin-go/docs"
	"github.com/mnmonherdene1234/files-gin-go/server"
)

// @title files-gin-go
// @version 1.0.0
// @description A simple file management API using Gin
func main() {
	server.Start()
}
