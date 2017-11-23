package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	SlackVerificationToken string `json:"slack_verification_token"`
	SlackAppClientId string `json:"slack_app_client_id"`
	SlackAppClientSecret string `json:"slack_app_client_secret"`
	DatabaseHost string `json:"database_host"`
	DatabasePort int `json:"database_port"`
	DatabaseUsername string `json:"database_username"`
	DatabasePassword string `json:"database_password"`
	DatabaseName string `json:"database_name"`
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
