package doc

import (
	"context"

	"github.com/Masterminds/squirrel"
	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
)

// GetListDocuments - получение списка документов для пользователя
func (r *Repository) GetListDocuments(ctx context.Context, userID string) ([]buisnesModel.Document, error) {
	query, args, err := r.sb.Select("*").From("documents").Where(squirrel.Eq{"user_id": userID}).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []buisnesModel.Document
	for rows.Next() {
		var doc buisnesModel.Document
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
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}
