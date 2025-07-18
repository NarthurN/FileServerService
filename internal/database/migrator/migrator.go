package migrator

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// Migrator handles database schema migrations.
type Migrator struct {
	db *sql.DB
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Up(ctx context.Context) error {
	// Установка провайдера для embedded файлов
	goose.SetBaseFS(embedMigrations)

	// Применение миграций
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("🚨 ошибка установки диалекта: %w", err)
	}

	if err := goose.UpContext(ctx, m.db, "migrations"); err != nil {
		return fmt.Errorf("🚨 ошибка применения миграций: %w", err)
	}

	log.Println("✅ Миграции применены успешно")
	return nil
}

func (m *Migrator) Down(ctx context.Context) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("🚨 ошибка установки диалекта: %w", err)
	}

	if err := goose.DownContext(ctx, m.db, "migrations"); err != nil {
		return fmt.Errorf("🚨 ошибка отмены миграций: %w", err)
	}

	log.Println("✅ Миграции отменены успешно")
	return nil
}

func (m *Migrator) Status(ctx context.Context) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("🚨 ошибка установки диалекта: %w", err)
	}

	if err := goose.StatusContext(ctx, m.db, "migrations"); err != nil {
		return fmt.Errorf("🚨 ошибка получения статуса миграций: %w", err)
	}

	return nil
}
