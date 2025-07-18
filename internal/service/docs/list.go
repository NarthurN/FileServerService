package docs

import (
	"context"

	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
)

// GetListDocuments - получение списка документов для пользователя
func (s *service) GetListDocuments(ctx context.Context, userID string) ([]buisnesModel.Document, error) {
	return s.repo.GetListDocuments(ctx, userID)
}
