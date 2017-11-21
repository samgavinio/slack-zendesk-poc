package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	SlackVerificationToken string `json:"slack_verification_token"`
}

func GetConfig() (config Config) {
	filePath, _ := filepath.Abs("config/parameters.json")
	file, fileError := os.Open(filePath)
	if fileError != nil {
		panic(fileError)
	}

	decoder := json.NewDecoder(file)
	jsonError := decoder.Decode(&config)
	if jsonError != nil {
		panic(jsonError)
	}

	return config
}
