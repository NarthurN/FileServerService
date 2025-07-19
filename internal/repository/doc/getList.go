package doc

import (
	"context"
	"log"

	"github.com/Masterminds/squirrel"
	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
)

// GetListDocuments - получение списка документов для пользователя
func (r *Repository) GetListDocuments(ctx context.Context, userID string) ([]buisnesModel.Document, error) {
	log.Printf("Repository: Получение документов для пользователя %s", userID)

	query, args, err := r.sb.Select("*").From("documents").Where(squirrel.Eq{"user_id": userID}).ToSql()
	if err != nil {
		log.Printf("Repository: Ошибка создания SQL запроса: %v", err)
		return nil, err
	}

	log.Printf("Repository: SQL запрос: %s", query)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("Repository: Ошибка выполнения SQL запроса: %v", err)
		return nil, err
	}
	defer rows.Close()

	var docs []buisnesModel.Document
	for rows.Next() {
		var doc buisnesModel.Document
		log.Printf("Repository: Сканирование документа...")
		if err := rows.Scan(
			&doc.ID,
			&doc.UserID,
			&doc.Name,
			&doc.MimeType,
			&doc.FilePath,
			&doc.IsFile,
			&doc.IsPublic,
			&doc.JSONData,
			&doc.Grants,
			&doc.CreatedAt,
			&doc.UpdatedAt,
		); err != nil {
			log.Printf("Repository: Ошибка сканирования строки: %v", err)
			return nil, err
		}
		log.Printf("Repository: Документ отсканирован: %s", doc.Name)
		docs = append(docs, doc)
	}

	log.Printf("Repository: Найдено %d документов", len(docs))
	return docs, nil
}
