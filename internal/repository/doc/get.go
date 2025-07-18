package doc

import (
	"context"

	"github.com/Masterminds/squirrel"
	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
	"github.com/jackc/pgx/v5"
)

// GetDocument - получение документа по ID
func (r *Repository) GetDocument(ctx context.Context, id string) (buisnesModel.Document, error) {
	query, args, err := r.sb.Select("*").From("documents").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return buisnesModel.Document{}, err
	}

	var doc buisnesModel.Document
	if err := r.pool.QueryRow(ctx, query, args...).Scan(
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
		if err == pgx.ErrNoRows {
			return buisnesModel.Document{}, buisnesModel.ErrNotFound
		}
		return buisnesModel.Document{}, err
	}

	return doc, nil
}
