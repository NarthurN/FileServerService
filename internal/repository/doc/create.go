package doc

import (
	"context"
	"log"

	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
)

// CreateDocument - создание документа
func (r *Repository) CreateDocument(ctx context.Context, doc buisnesModel.Document) (buisnesModel.Document, error) {
	log.Printf("RepLayer: Начало загрузки документа %s\n", doc.Name)
	query, args, err := r.sb.Insert("documents").
		Columns("id", "user_id", "name", "mime_type", "file_path", "is_file", "is_public", "json_data", "grants", "created_at", "updated_at").
		Values(doc.ID, doc.UserID, doc.Name, doc.MimeType, doc.FilePath, doc.IsFile, doc.IsPublic, doc.JSONData, doc.Grants, doc.CreatedAt, doc.UpdatedAt).
		ToSql()
	if err != nil {
		log.Printf("RepLayer: ошибка подготовки запроса загрзки документа%s: %v \n", doc.Name, err)
		return buisnesModel.Document{}, err
	}
	log.Printf("RepLayer: Запрос для загрузки документа %s подготовлен \n", doc.Name)
	if _, err := r.pool.Exec(ctx, query, args...); err != nil {
		log.Printf("RepLayer: ошибка загрузки документа %s: %v \n", doc.Name, err)
		return buisnesModel.Document{}, err
	}
	log.Printf("RepLayer: Документ %s загружен \n", doc.Name)
	return doc, nil
}
