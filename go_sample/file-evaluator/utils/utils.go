package utils

import (
	"encoding/json"
	"log"
	"os"
)

func LoadConfig[T any](path string, cfg *T) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read config file: %s", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		log.Fatalf("failed to parse config file: %s", err)
	}
}
