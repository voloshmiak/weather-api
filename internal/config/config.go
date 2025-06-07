package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"path/filepath"
)

type Database struct {
	User       string `env:"DB_USER" env-default:"postgres"`
	Password   string `env:"DB_PASSWORD" env-default:"your_password_here"`
	Host       string `env:"DB_HOST" env-default:"localhost"`
	Port       string `env:"DB_PORT" env-default:"5432"`
	Name       string `env:"DB_NAME" env-default:"forum_database"`
	Migrations string `env:"MIGRATIONS_PATH" env-default:"./internal/database/migrations"`
}

func (d *Database) MigrationURL() string {
	migrationsAbsPath, _ := filepath.Abs(d.Migrations)
	migrationsSlashPath := filepath.ToSlash(migrationsAbsPath)
	return fmt.Sprintf("file://%s", migrationsSlashPath)
}

type Config struct {
	DB     Database
	Server struct {
		Port         string `env:"API_PORT" env-default:"8080"`
		ReadTimeout  int    `env:"SERVER_READ_TIMEOUT" env-default:"5"`
		WriteTimeout int    `env:"SERVER_WRITE_TIMEOUT" env-default:"10"`
		IdleTimeout  int    `env:"SERVER_IDLE_TIMEOUT" env-default:"15"`
	}
	Mail struct {
		Host string `env:"SMTP_HOST" env-default:"mailhog"`
		Port string `env:"SMTP_PORT" env-default:"1025"`
	}
	WeatherAPIKey string `env:"WEATHER_API_KEY" env-default:"your_weather_api_key_here"`
}

func New() (*Config, error) {
	config := new(Config)

	if err := cleanenv.ReadEnv(config); err != nil {
		return nil, err
	}

	return config, nil
}
