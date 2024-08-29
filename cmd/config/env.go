package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string

	DbUser                 string
	DbPassword             string
	DbAddress              string
	DbName                 string
	JWTExpirationInSeconds int64
	JWTSecret              string
}

// create a singleton
var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "8080"),
		DbUser:                 getEnv("DB_USER", "root"),
		DbPassword:             getEnv("DB_PASSWORD", "***"),
		DbAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DbName:                 getEnv("DB_NAME", "go-test-db"),
		JWTExpirationInSeconds: getEnvInt("JWT_EXP", 3600*24*7),
		JWTSecret:              getEnv("JWT_SECRET", "secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return intValue
	}

	return fallback
}
