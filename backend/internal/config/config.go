package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() Config {
	c := Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://codenames:codenames@localhost:5432/codenames?sslmode=disable"),
	}
	return c
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
