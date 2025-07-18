package repository

import (
	"context"

	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
)

type FileServerRepository interface {
	CreateDocument(ctx context.Context, doc buisnesModel.Document) (buisnesModel.Document, error)
	GetDocument(ctx context.Context, id string) (buisnesModel.Document, error)
	GetListDocuments(ctx context.Context, userID string) ([]buisnesModel.Document, error)
}
