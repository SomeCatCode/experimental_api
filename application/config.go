package application

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	MongoUri string
	MongoDb  string
}

func loadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Laden der Umgebungsvariablen: %w", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if port[0] != ':' {
		port = ":" + port
	}

	mongoUri := os.Getenv("MONGO_URI")
	if mongoUri == "" {
		mongoUri = "mongodb://localhost:27017"
	}

	mongoDb := os.Getenv("MONGO_DB")
	if mongoDb == "" {
		mongoDb = "experimental_api"
	}

	return &Config{
		Port:     port,
		MongoUri: mongoUri,
		MongoDb:  mongoDb,
	}, nil
}
