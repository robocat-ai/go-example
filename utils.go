package main

import (
	"log"
	"os"
	"path"
	"path/filepath"

	robocat "github.com/robocat-ai/robocat/pkg/client"
)

func outputPath(path string) string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %s", err)
	}

	return filepath.Join(pwd, "output", path)
}

func writeFile(file *robocat.File) {

	destinationPath := outputPath(file.Path)

	err := os.MkdirAll(path.Dir(destinationPath), 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(destinationPath, file.Payload, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
