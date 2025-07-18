package doc

import (
	"context"
	"log"

	"github.com/Masterminds/squirrel"
	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
)

// DeleteDocument - удаление документа по ID
func (r *Repository) DeleteDocument(ctx context.Context, id string) error {
	log.Printf("RepLayer: Начало удаления документа %s\n", id)

	query, args, err := r.sb.Delete("documents").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		log.Printf("RepLayer: ошибка подготовки запроса удаления документа %s: %v\n", id, err)
		return err
	}

	log.Printf("RepLayer: Запрос для удаления документа %s подготовлен\n", id)

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("RepLayer: ошибка удаления документа %s: %v\n", id, err)
		return err
	}

	// Проверяем, что документ был удален
	if result.RowsAffected() == 0 {
		log.Printf("RepLayer: документ %s не найден для удаления\n", id)
		return buisnesModel.ErrNotFound
	}

	log.Printf("RepLayer: Документ %s удален\n", id)
	return nil
}
