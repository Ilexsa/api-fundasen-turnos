package config

import (
	"strconv"
)

type AppConfig struct {
	APIAddr      string
	APIToken     string
	AllowedOrigins string

	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBEncrypt string
	DBConnTimeoutSeconds int
}

func Load() AppConfig {
	return AppConfig{
		APIAddr:        Get("API_ADDR", ":8080"),
		APIToken:       Get("API_TOKEN", ""),
		AllowedOrigins: Get("ALLOWED_ORIGINS", "*"),

		DBHost:    Get("DB_HOST", Get("DB_HOST", "localhost")),
		DBPort:    Get("DB_PORT", Get("DB_PORT", "1433")),
		DBUser:    Get("DB_USER", Get("DB_USER", "sa")),
		DBPass:    Get("DB_PASS", Get("DB_PASS", "")),
		DBName:    Get("DB_NAME", Get("DB_NAME", "TURNOS")),
		DBEncrypt: Get("DB_ENCRYPT", Get("DB_ENCRYPT", "disable")),
		DBConnTimeoutSeconds: getInt("DB_CONN_TIMEOUT_SECONDS", 10),
	}
}

func getInt(key string, def int) int {
	v := Get(key, "")
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
