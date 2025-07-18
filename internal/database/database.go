package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/NarthurN/FileServerService/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
