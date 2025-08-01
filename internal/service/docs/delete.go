package docs

import (
	"context"
	"fmt"
	"log"
)

// DeleteDocument - удаление документа с проверкой прав
func (s *service) DeleteDocument(ctx context.Context, id string) error {
	log.Printf("ServiceLayer: Удаление документа %s", id)

	if id == "" {
		return fmt.Errorf("document ID is required")
	}

	// Получаем документ для проверки существования
	doc, err := s.repo.GetDocument(ctx, id)
	if err != nil {
		log.Printf("ServiceLayer: Документ %s не найден для удаления: %v", id, err)
		return fmt.Errorf("document not found: %w", err)
	}

	// Удаляем документ
	if err := s.repo.DeleteDocument(ctx, id); err != nil {
		log.Printf("ServiceLayer: Ошибка удаления документа %s: %v", id, err)
		return fmt.Errorf("failed to delete document: %w", err)
	}

	// Инвалидируем кэш для документа и пользователя
	if err := s.cacheManager.InvalidateDocument(ctx, id); err != nil {
		log.Printf("ServiceLayer: Ошибка инвалидации кэша документа: %v", err)
	}
	if err := s.cacheManager.InvalidateUserDocuments(ctx, doc.UserID); err != nil {
		log.Printf("ServiceLayer: Ошибка инвалидации кэша документов пользователя: %v", err)
	}

	log.Printf("ServiceLayer: Документ %s (%s) успешно удален, кэш инвалидирован", doc.Name, id)
	return nil
}
