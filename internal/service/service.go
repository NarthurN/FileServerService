package service

import (
	"context"

	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
)

// FileServerService - интерфейс сервиса для работы с документами
type FileServerService interface {
	CreateDocument(ctx context.Context, doc buisnesModel.Document) (buisnesModel.Document, error)
	GetDocument(ctx context.Context, id string) (buisnesModel.Document, error)
	GetListDocuments(ctx context.Context, userID string) ([]buisnesModel.Document, error)
}
