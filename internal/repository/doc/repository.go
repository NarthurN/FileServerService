package doc

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Masterminds/squirrel"
)

// Repository - репозиторий для работы с документами
type Repository struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
}

// NewRepository - создание нового репозитория
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
