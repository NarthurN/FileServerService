package docs

import (
	"context"

	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
)

// CreateDocument - создание документа
func (s *service) CreateDocument(ctx context.Context, doc buisnesModel.Document) (buisnesModel.Document, error) {
	return s.repo.CreateDocument(ctx, doc)
}
