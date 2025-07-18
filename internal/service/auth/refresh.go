package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/NarthurN/FileServerService/internal/model"
	"github.com/google/uuid"
)

// RefreshToken - обновление токена
func (s *Service) RefreshToken(ctx context.Context, oldToken string) (string, error) {
	log.Printf("AuthService: Обновление токена")

	// Валидируем старый токен
	user, err := s.ValidateToken(ctx, oldToken)
	if err != nil {
		return "", fmt.Errorf("invalid old token: %w", err)
	}

	// Деактивируем старый токен
	if err := s.repo.DeactivateToken(ctx, oldToken); err != nil {
		log.Printf("AuthService: Предупреждение - не удалось деактивировать старый токен: %v", err)
	}

	// Создаем новый токен
	newTokenValue, err := s.generateSecureToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate new token: %w", err)
	}

	newToken := model.Token{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     newTokenValue,
		ExpiresAt: time.Now().UTC().Add(s.getTokenLifetime()),
		CreatedAt: time.Now().UTC(),
		IsActive:  true,
	}

	if _, err := s.repo.CreateToken(ctx, newToken); err != nil {
		return "", fmt.Errorf("failed to save new token: %w", err)
	}

	log.Printf("AuthService: Токен успешно обновлен")
	return newTokenValue, nil
}
