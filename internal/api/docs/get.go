package docs

import (
	"context"

	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
)

// GetDocument - получение документа по ID
func (a *api) GetDocument(ctx context.Context, id string) (buisnesModel.Document, error) {
	return a.docsService.GetDocument(ctx, id)
}	
