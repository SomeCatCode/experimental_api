package application

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     uint16
	MongoUri string
	MongoDb  string
}

func LoadConfig() Config {
	_ = godotenv.Load()

	cfg := Config{
		Port:     8080,
		MongoUri: "mongodb://localhost:27017",
		MongoDb:  "experimental_api",
	}

	if serverPort, exists := os.LookupEnv("PORT"); exists {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.Port = uint16(port)
		}
	}

	if mongoUri, exists := os.LookupEnv("MONGO_URI"); exists {
		cfg.MongoUri = mongoUri
	}

	if mongoDb, exists := os.LookupEnv("MONGO_DB"); exists {
		cfg.MongoDb = mongoDb
	}

	return cfg
}
