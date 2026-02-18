package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Port      string
	MongoURI  string
	MongoDB   string
	JWTSecret string
}

func LoadEnv() Env {
	_ = godotenv.Load()

	env := Env{
		Port:      getEnv("PORT", "8080"),
		MongoURI:  os.Getenv("MONGO_URI"),
		MongoDB:   getEnv("MONGO_DB", "cardflex"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	if env.MongoURI == "" {
		log.Fatal("MONGO_URI is required")
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
