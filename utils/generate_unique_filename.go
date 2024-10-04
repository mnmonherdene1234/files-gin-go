package utils

import (
	"fmt"
	"path/filepath"
	"time"
)

func GenerateUniqueFilename(originalFilename string) string {
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(originalFilename)
	name := filepath.Base(originalFilename[:len(originalFilename)-len(ext)])
	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}
