package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/NarthurN/FileServerService/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// NewPool создает пул соединений к PostgreSQL
func NewPool(cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("🚨 ошибка парсинга строки подключения: %w", err)
	}

	poolConfig.MaxConns = 30
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = time.Minute * 30
	poolConfig.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("🚨 ошибка создания пула соединений: %w", err)
	}

	// Проверка соединения
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("🚨 ошибка пинга базы данных: %w", err)
	}

	log.Println("🚀 Успешное подключение к PostgreSQL")
	return pool, nil
}

// NewSQLDB создает *sql.DB для миграций Goose
func NewSQLDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("🚨 ошибка открытия базы данных: %w", err)
	}

	// Настройка соединения
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("🚨 ошибка пинга базы данных: %w", err)
	}

	log.Println("✅ SQL соединение для миграций создано")
	return db, nil
}

// ConvertPoolToSQLDB конвертирует pgx.Pool в *sql.DB (если нужно)
func ConvertPoolToSQLDB(pool *pgxpool.Pool) *sql.DB {
	return stdlib.OpenDB(*pool.Config().ConnConfig)
}
