package auth

import (
	"context"
	"fmt"
	"log"
)

// LogoutUser - завершение сессии
func (s *Service) LogoutUser(ctx context.Context, tokenValue string) error {
	log.Printf("AuthService: Завершение сессии для токена")

	if tokenValue == "" {
		return fmt.Errorf("token is required")
	}

	// Проверяем, что токен существует
	_, err := s.repo.GetTokenByValue(ctx, tokenValue)
	if err != nil {
		log.Printf("AuthService: Токен для выхода не найден: %v", err)
		return fmt.Errorf("token not found")
	}

	// Деактивируем токен
	if err := s.repo.DeactivateToken(ctx, tokenValue); err != nil {
		log.Printf("AuthService: Ошибка деактивации токена: %v", err)
		return fmt.Errorf("failed to deactivate token: %w", err)
	}

	log.Printf("AuthService: Сессия успешно завершена")
	return nil
}
