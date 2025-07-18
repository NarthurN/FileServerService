package auth

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/NarthurN/FileServerService/internal/model"
	"github.com/google/uuid"
)

// AuthenticateUser - аутентификация с полной проверкой
func (s *Service) AuthenticateUser(ctx context.Context, login, password string) (string, error) {
	log.Printf("AuthService: Начало аутентификации пользователя %s", login)

	// Валидация входных данных
	if err := s.validateAuthInput(login, password); err != nil {
		log.Printf("AuthService: Неверные входные данные для аутентификации: %v", err)
		return "", fmt.Errorf("invalid input: %w", err)
	}

	// Нормализация логина
	normalizedLogin := strings.ToLower(strings.TrimSpace(login))

	// Получение пользователя
	user, err := s.repo.GetUserByLogin(ctx, normalizedLogin)
	if err != nil {
		log.Printf("AuthService: Пользователь %s не найден: %v", normalizedLogin, err)
		return "", fmt.Errorf("invalid credentials")
	}

	// Проверка пароля
	if err := s.verifyPassword(password, user.Password); err != nil {
		log.Printf("AuthService: Неверный пароль для пользователя %s", normalizedLogin)
		return "", fmt.Errorf("invalid credentials")
	}

	// Бизнес-логика: деактивация старых токенов (опционально)
	if err := s.cleanupOldTokens(ctx, user.ID); err != nil {
		log.Printf("AuthService: Предупреждение - не удалось очистить старые токены: %v", err)
		// Не прерываем процесс, это не критичная ошибка
	}

	// Генерация нового токена
	tokenValue, err := s.generateSecureToken()
	if err != nil {
		log.Printf("AuthService: Ошибка генерации токена: %v", err)
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Создание токена с настраиваемым временем жизни
	token := model.Token{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     tokenValue,
		ExpiresAt: time.Now().UTC().Add(s.getTokenLifetime()),
		CreatedAt: time.Now().UTC(),
		IsActive:  true,
	}

	if _, err := s.repo.CreateToken(ctx, token); err != nil {
		log.Printf("AuthService: Ошибка сохранения токена: %v", err)
		return "", fmt.Errorf("failed to save token: %w", err)
	}

	log.Printf("AuthService: Пользователь %s успешно аутентифицирован", normalizedLogin)
	return tokenValue, nil
}
