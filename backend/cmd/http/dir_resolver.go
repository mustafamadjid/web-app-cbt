package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ResolveAppDir() string {
	appDir := strings.TrimSpace(os.Getenv("APP_DIR"))
	if appDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("resolve APP_DIR from working dir: %v", err)
		}
		appDir = wd
	}

	absAppDir, err := filepath.Abs(appDir)
	if err != nil {
		log.Fatalf("resolve APP_DIR to absolute path: %v", err)
	}

	return filepath.Clean(absAppDir)
}

func ResolveUploadDir(appDir string) string {
	uploadDir := strings.TrimSpace(os.Getenv("UPLOAD_DIR"))
	if uploadDir == "" {
		return filepath.Join(appDir, "public", "uploads")
	}

	if filepath.IsAbs(uploadDir) {
		return filepath.Clean(uploadDir)
	}

	return filepath.Join(appDir, uploadDir)
}
