package docs

import (
	"context"
	"fmt"
	"log"

	"github.com/NarthurN/FileServerService/internal/model"
)

// GetDocument - получение документа с проверкой прав доступа
func (s *service) GetDocument(ctx context.Context, id string) (model.Document, error) {
	log.Printf("ServiceLayer: Получение документа %s", id)

	if id == "" {
		return model.Document{}, fmt.Errorf("document ID is required")
	}

	doc, err := s.repo.GetDocument(ctx, id)
	if err != nil {
		log.Printf("ServiceLayer: Документ %s не найден: %v", id, err)
		return model.Document{}, fmt.Errorf("document not found: %w", err)
	}

	log.Printf("ServiceLayer: Документ %s найден", id)
	return doc, nil
}

// GetDocumentsForUser - получение документов с учетом прав доступа
func (s *service) GetDocumentsForUser(ctx context.Context, requestUserID, targetUserID string) ([]model.Document, error) {
	log.Printf("ServiceLayer: Получение документов пользователя %s для пользователя %s", targetUserID, requestUserID)

	// Если запрашивает свои документы
	if requestUserID == targetUserID {
		return s.GetListDocuments(ctx, targetUserID)
	}

	// Получаем все документы целевого пользователя
	allDocs, err := s.repo.GetListDocuments(ctx, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get target user documents: %w", err)
	}

	// Получаем пользователя-запросчика для проверки его логина
	requestUser, err := s.repo.GetUserByID(ctx, requestUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get request user: %w", err)
	}

	// Фильтруем документы: публичные + те, к которым есть доступ
	var accessibleDocs []model.Document
	for _, doc := range allDocs {
		if s.hasAccessToDocument(doc, requestUser.Login) {
			accessibleDocs = append(accessibleDocs, doc)
		}
	}

	log.Printf("ServiceLayer: Пользователь %s имеет доступ к %d из %d документов", requestUserID, len(accessibleDocs), len(allDocs))
	return s.sortDocuments(accessibleDocs), nil
}
