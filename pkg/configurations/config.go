package configurations

import (
	"os"
	"log"
	"encoding/json"
)

type Config struct {
	Port			string `json:"port"`
	AllowedOrigins	string `json:"allowedOrigins"`
	DbURL			string `json:"dbURL"`
	DbName			string `json:"dbName"`
	CollectionName	string `json:"collectionName"`
	AuthUser		string `json:"authUser"`
	AuthPassword	string `json:"authPassword"`
}

type Secrets struct {
	AuthUser		string `json:"authUser"`
	AuthPassword	string `json:"authPassword"`
	DbURL			string `json:"dbURL"`
}

func LoadConfig() (config *Config) {
	configPath := "./config.json"
	secretsPath := "./secrets.json"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Failed to find config file at path \"%v\" (Stat Error = \"%v\".\n", configPath, err.Error())
	}
	if _, err := os.Stat(secretsPath); os.IsNotExist(err) {
		log.Fatalf("Failed to find secrets file at path \"%v\" (Stat Error = \"%v\".\n", secretsPath, err.Error())
	}
	confFile, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Error opening config file: \n\"%v\".\n", err.Error())
	}
	defer confFile.Close()
	secretsFile, err := os.Open(secretsPath)
	if err != nil {
		log.Fatalf("Error opening secrets file: \n\"%v\".\n", err.Error())
	}
	defer secretsFile.Close()
	if err = json.NewDecoder(confFile).Decode(&config); err != nil {
		log.Fatalf("Error decoding config file: \n\"%v\".\n", err.Error())
	}
	var secrets Secrets
	if err = json.NewDecoder(secretsFile).Decode(&secrets); err != nil {
		log.Fatalf("Error decoding secrets file: \n\"%v\".\n", err.Error())
	}
	log.Print("Secrets successfully loaded.")
	config.DbURL = secrets.DbURL
	config.AuthUser = secrets.AuthUser
	config.AuthPassword = secrets.AuthPassword
	log.Print("Config successfully loaded.")
	return config
}