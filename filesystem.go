package main

import (
	"fmt"
	"os"
)

const filePath = "files/"

func getFile() (string, error) {
	filepath := filePath + "test.txt"
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("erreur de lecture du fichier: %w", err)
	}
	return string(data), nil
}
