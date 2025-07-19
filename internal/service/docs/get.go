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

	// Пытаемся получить из кэша
	if cachedDoc, found := s.cacheManager.GetDocument(ctx, id); found {
		if doc, ok := cachedDoc.(model.Document); ok {
			log.Printf("ServiceLayer: Документ %s найден в кэше", id)
			return doc, nil
		}
	}

	// Если нет в кэше, получаем из БД
	doc, err := s.repo.GetDocument(ctx, id)
	if err != nil {
		log.Printf("ServiceLayer: Документ %s не найден: %v", id, err)
		return model.Document{}, fmt.Errorf("document not found: %w", err)
	}

	// Сохраняем в кэш
	if err := s.cacheManager.SetDocument(ctx, id, doc); err != nil {
		log.Printf("ServiceLayer: Ошибка сохранения документа %s в кэш: %v", id, err)
	}

	log.Printf("ServiceLayer: Документ %s найден и сохранен в кэш", id)
	return doc, nil
}

// GetDocumentsForUser - получение документов с учетом прав доступа
func (s *service) GetDocumentsForUser(ctx context.Context, requestUserID, targetUserID string) ([]model.Document, error) {
	log.Printf("ServiceLayer: Получение документов пользователя %s для пользователя %s", targetUserID, requestUserID)

	// Если запрашивает свои документы
	if requestUserID == targetUserID {
		log.Printf("ServiceLayer: Запрос собственных документов")
		return s.GetListDocuments(ctx, targetUserID)
	}

	// Получаем все документы целевого пользователя
	log.Printf("ServiceLayer: Получение документов целевого пользователя %s", targetUserID)
	allDocs, err := s.repo.GetListDocuments(ctx, targetUserID)
	if err != nil {
		log.Printf("ServiceLayer: Ошибка получения документов целевого пользователя: %v", err)
		return nil, fmt.Errorf("failed to get target user documents: %w", err)
	}
	log.Printf("ServiceLayer: Найдено %d документов целевого пользователя", len(allDocs))

	// Получаем пользователя-запросчика для проверки его логина
	log.Printf("ServiceLayer: Получение пользователя-запросчика %s", requestUserID)
	requestUser, err := s.repo.GetUserByID(ctx, requestUserID)
	if err != nil {
		log.Printf("ServiceLayer: Ошибка получения пользователя-запросчика: %v", err)
		return nil, fmt.Errorf("failed to get request user: %w", err)
	}
	log.Printf("ServiceLayer: Пользователь-запросчик найден: %s", requestUser.Login)

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
