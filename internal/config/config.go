package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Printf("Warning: Environment variable %s not set, using fallback '%s'", key, fallback)
	return fallback
}

func GetDatabaseURL() string {
	dbHost := getEnv("DB_HOST", "db")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "user")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "weather_database")

	encodedPassword := url.QueryEscape(dbPassword)

	dbURLForMigrate := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, encodedPassword, dbHost, dbPort, dbName)

	return dbURLForMigrate
}

func GetMigrationURL() string {
	migrationsPath, _ := filepath.Abs("./migrations")
	migrationsPathURL := filepath.ToSlash(migrationsPath)
	migrationURL := "file:///" + migrationsPathURL

	return migrationURL
}

func GetAPIPort() string {
	apiPort := getEnv("API_PORT", "8080")
	return apiPort
}

func GetSMTPHost() string {
	smtpHost := getEnv("SMTP_HOST", "mailhog")
	return smtpHost
}

func GetSMTPPort() string {
	smtpPort := getEnv("SMTP_PORT", "1025")
	return smtpPort
}
