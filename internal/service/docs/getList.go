package docs

import (
	"context"
	"fmt"
	"log"

	"github.com/NarthurN/FileServerService/internal/model"
)

// GetListDocuments - получение списка документов с сортировкой
func (s *service) GetListDocuments(ctx context.Context, userID string) ([]model.Document, error) {
	log.Printf("ServiceLayer: Получение списка документов для пользователя %s", userID)

	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Пытаемся получить из кэша
	if cachedDocs, found := s.cacheManager.GetDocumentList(ctx, userID, "", "", 0); found {
		docs := make([]model.Document, len(cachedDocs))
		for i, docInterface := range cachedDocs {
			if doc, ok := docInterface.(model.Document); ok {
				docs[i] = doc
			}
		}
		log.Printf("ServiceLayer: Список документов для пользователя %s найден в кэше", userID)
		return docs, nil
	}

	// Проверяем, что пользователь существует
	if _, err := s.repo.GetUserByID(ctx, userID); err != nil {
		log.Printf("ServiceLayer: Пользователь %s не найден: %v", userID, err)
		return nil, fmt.Errorf("user not found: %w", err)
	}

	docs, err := s.repo.GetListDocuments(ctx, userID)
	if err != nil {
		log.Printf("ServiceLayer: Ошибка получения документов: %v", err)
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}

	// Применяем бизнес-логику сортировки (согласно заданию - по имени и дате создания)
	sortedDocs := s.sortDocuments(docs)

	// Сохраняем в кэш
	docsInterface := make([]interface{}, len(sortedDocs))
	for i, doc := range sortedDocs {
		docsInterface[i] = doc
	}
	if err := s.cacheManager.SetDocumentList(ctx, userID, "", "", 0, docsInterface); err != nil {
		log.Printf("ServiceLayer: Ошибка сохранения списка документов в кэш: %v", err)
	}

	log.Printf("ServiceLayer: Найдено %d документов для пользователя %s и сохранено в кэш", len(sortedDocs), userID)
	return sortedDocs, nil
}
