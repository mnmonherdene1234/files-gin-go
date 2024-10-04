package main

import (
	_ "github.com/mnmonherdene1234/files-gin-go/docs"
	"github.com/mnmonherdene1234/files-gin-go/server"
)

// @title GIN Files API
// @version 1.0.0
// @description A files server
func main() {
	server.Start()
}
