package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Printf("Warning: Environment variable %s not set, using fallback '%s'", key, fallback)
	return fallback
}

func NewPostgresDB() (*sql.DB, error) {
	dbHost := getEnv("DB_HOST", "db")                   // Default to 'db' for Docker Compose service name
	dbPort := getEnv("DB_PORT", "5432")                 // Default PostgreSQL port
	dbUser := getEnv("DB_USER", "youruser")             // Ensure this matches docker-compose.yml or your .env
	dbPassword := getEnv("DB_PASSWORD", "yourpassword") // Ensure this matches docker-compose.yml or your .env
	dbName := getEnv("DB_NAME", "weather_database")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	conn, err := sql.Open("pgx", psqlInfo)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // 15-second timeout
	defer cancel()

	err = conn.PingContext(ctx) // Use PingContext
	if err != nil {
		// This will now explicitly show a timeout error or other connection errors
		log.Fatalf("FATAL: Error pinging database: %v. Connection string was: %s", err, psqlInfo)
	}

	return conn, nil
}
