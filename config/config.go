package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	SlackAppClientId string `json:"slack_app_client_id"`
	SlackAppClientSecret string `json:"slack_app_client_secret"`
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
