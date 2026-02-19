package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Port      string
	DBDriver  string
	DBDSN     string
	JWTSecret string
}

func LoadEnv() Env {
	_ = godotenv.Load()

	env := Env{
		Port:      getEnv("PORT", "8080"),
		DBDriver:  getEnv("DB_DRIVER", "sqlite"),
		DBDSN:     getEnv("DB_DSN", "cardflex.db"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	if env.DBDSN == "" {
		log.Fatal("DB_DSN is required")
	}

	if env.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return env
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
