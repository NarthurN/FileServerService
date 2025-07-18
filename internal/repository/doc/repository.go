package doc

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Masterminds/squirrel"

	def "github.com/NarthurN/FileServerService/internal/repository"
)

// Проверка на реализацию интерфейса
var _ def.FileServerRepository = (*Repository)(nil)

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
