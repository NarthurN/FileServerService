package docs

import (
	"context"
	"fmt"
)

// HasAccessToDocument - проверка прав доступа к документу
func (s *service) HasAccessToDocument(ctx context.Context, userID, documentID string) (bool, error) {
	// Получаем документ
	doc, err := s.repo.GetDocument(ctx, documentID)
	if err != nil {
		return false, fmt.Errorf("document not found: %w", err)
	}

	// Владелец всегда имеет доступ
	if doc.UserID == userID {
		return true, nil
	}

	// Публичные документы доступны всем
	if doc.IsPublic {
		return true, nil
	}

	// Проверяем grants
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("user not found: %w", err)
	}

	return s.hasAccessToDocument(doc, user.Login), nil
}
