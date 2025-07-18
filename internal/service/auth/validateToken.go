package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/NarthurN/FileServerService/internal/model"
)

// ValidateToken - валидация токена с проверкой срока действия
func (s *Service) ValidateToken(ctx context.Context, tokenValue string) (model.User, error) {
	if tokenValue == "" {
		return model.User{}, fmt.Errorf("token is required")
	}

	// Получение токена с проверкой активности и срока действия
	token, err := s.repo.GetTokenByValue(ctx, tokenValue)
	if err != nil {
		log.Printf("AuthService: Токен не найден или недействителен: %v", err)
		return model.User{}, fmt.Errorf("invalid token")
	}

	// Дополнительная проверка срока действия (на всякий случай)
	if time.Now().UTC().After(token.ExpiresAt) {
		log.Printf("AuthService: Токен истек: %v", token.ExpiresAt)
		// Деактивируем истекший токен
		_ = s.repo.DeactivateToken(ctx, tokenValue)
		return model.User{}, fmt.Errorf("token expired")
	}

	// Получение пользователя
	user, err := s.repo.GetUserByID(ctx, token.UserID)
	if err != nil {
		log.Printf("AuthService: Пользователь для токена не найден: %v", err)
		return model.User{}, fmt.Errorf("user not found")
	}

	return user, nil
}
