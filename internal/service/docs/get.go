package docs

import (
	"context"

	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
)

// GetDocument - получение документа по ID
func (s *service) GetDocument(ctx context.Context, id string) (buisnesModel.Document, error) {
	return s.repo.GetDocument(ctx, id)
}
