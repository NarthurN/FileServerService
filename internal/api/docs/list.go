package docs

import (
	"context"

	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
)

// GetListDocuments - получение списка документов для пользователя
func (a *api) GetListDocuments(ctx context.Context, userID string) ([]buisnesModel.Document, error) {
	return a.docsService.GetListDocuments(ctx, userID)
}
