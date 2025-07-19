package docs

import (
	"context"
	"fmt"
)

// HasAccessToDocument - проверка прав доступа к документу
func (s *service) HasAccessToDocument(ctx context.Context, userID, documentID string) (bool, error) {
	// Пытаемся получить из кэша
	if hasAccess, found := s.cacheManager.GetAccess(ctx, documentID, userID); found {
		return hasAccess, nil
	}

	// Получаем документ
	doc, err := s.repo.GetDocument(ctx, documentID)
	if err != nil {
		return false, fmt.Errorf("document not found: %w", err)
	}

	// Владелец всегда имеет доступ
	if doc.UserID == userID {
		// Сохраняем в кэш
		s.cacheManager.SetAccess(ctx, documentID, userID, true)
		return true, nil
	}

	// Публичные документы доступны всем
	if doc.IsPublic {
		// Сохраняем в кэш
		s.cacheManager.SetAccess(ctx, documentID, userID, true)
		return true, nil
	}

	// Проверяем grants
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("user not found: %w", err)
	}

	hasAccess := s.hasAccessToDocument(doc, user.Login)

	// Сохраняем в кэш
	s.cacheManager.SetAccess(ctx, documentID, userID, hasAccess)

	return hasAccess, nil
}
