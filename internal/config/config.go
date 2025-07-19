package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Все настройки
type Config struct {
	Database DatabaseConfig // База данных
	Server   ServerConfig   // Сервер
	Auth     AuthConfig     // Авторизация админа
}

// Настройки базы данных
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Настройки сервера
type ServerConfig struct {
	Host string
	Port int
}

// Настройки авторизации для админа
type AuthConfig struct {
	AdminToken    string        // Фиксированный токен администратора для регистрации
	TokenLifetime time.Duration // Время жизни пользовательских токенов
	JWTSecret     string        // Секрет для JWT
}

func Load() (*Config, error) {
	// Пытаемся загрузить .env файл, но не возвращаем ошибку если его нет
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️ Файл .env не найден, используем переменные окружения:", err)
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "docs_server"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnvInt("SERVER_PORT", 8080),
		},
		Auth: AuthConfig{
			AdminToken:    getEnv("ADMIN_TOKEN", "admin-secret-token-123456"),
			TokenLifetime: getTokenLifetime(),
			JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getTokenLifetime() time.Duration {
	// По умолчанию 24 часа
	defaultHours := 24

	if hoursStr := os.Getenv("TOKEN_LIFETIME_HOURS"); hoursStr != "" {
		if hours, err := strconv.Atoi(hoursStr); err == nil && hours > 0 {
			return time.Duration(hours) * time.Hour
		}
	}

	return time.Duration(defaultHours) * time.Hour
}
