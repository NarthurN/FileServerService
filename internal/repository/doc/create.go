package doc

import (
	"context"

	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
)

// CreateDocument - создание документа
func (r *Repository) CreateDocument(ctx context.Context, doc buisnesModel.Document) (buisnesModel.Document, error) {
	query, args, err := r.sb.Insert("documents").
		Columns("id", "user_id", "name", "mime_type", "file_path", "is_file", "is_public", "json_data", "grants", "created_at", "updated_at").
		Values(doc.ID, doc.UserID, doc.Name, doc.MimeType, doc.FilePath, doc.IsFile, doc.IsPublic, doc.JSONData, doc.Grants, doc.CreatedAt, doc.UpdatedAt).
		ToSql()
	if err != nil {
		return buisnesModel.Document{}, err
	}

	if _, err := r.pool.Exec(ctx, query, args...); err != nil {
		return buisnesModel.Document{}, err
	}
	return doc, nil
}
