package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GatewayPort    string
	JWTSecret      string
	AuthAPIURL     string
	CryptoETLURL   string
	InternalAPIKey string
}

var Cfg Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error cargando .env")
	}

	Cfg = Config{
		GatewayPort:    os.Getenv("GATEWAY_PORT"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		AuthAPIURL:     os.Getenv("AUTH_API_URL"),
		CryptoETLURL:   os.Getenv("CRYPTO_ETL_URL"),
		InternalAPIKey: os.Getenv("INTERNAL_API_KEY"),
	}
}
