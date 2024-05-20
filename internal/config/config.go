package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct holds the configuration variables
type Config struct {
	PrivateKey string
	RPCUrl     string
}

// LoadConfig loads the configuration from .env file
func LoadConfig() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Return config struct with loaded values
	return &Config{
		PrivateKey: os.Getenv("PRIVATE_KEY"),
		RPCUrl:     os.Getenv("RPC_URL"),
	}
}
