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

// RegisterUser - регистрация нового пользователя с полной валидацией
func (s *Service) RegisterUser(ctx context.Context, adminToken, login, password string) (model.User, error) {
	log.Printf("AuthService: Начало регистрации пользователя %s", login)

	// Бизнес-валидация: проверка админского токена
	if err := s.validateAdminToken(adminToken); err != nil {
		log.Printf("AuthService: Неверный админский токен")
		return model.User{}, fmt.Errorf("invalid admin token: %w", err)
	}

	// Валидация логина согласно заданию
	if err := s.validateLogin(login); err != nil {
		log.Printf("AuthService: Неверный логин %s: %v", login, err)
		return model.User{}, fmt.Errorf("invalid login: %w", err)
	}

	// Валидация пароля согласно заданию
	if err := s.validatePassword(password); err != nil {
		log.Printf("AuthService: Неверный пароль для пользователя %s: %v", login, err)
		return model.User{}, fmt.Errorf("invalid password: %w", err)
	}

	// Бизнес-правило: проверяем уникальность логина
	if err := s.checkLoginUniqueness(ctx, login); err != nil {
		log.Printf("AuthService: Логин %s уже занят", login)
		return model.User{}, fmt.Errorf("login already taken: %w", err)
	}

	// Бизнес-логика: хеширование пароля
	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		log.Printf("AuthService: Ошибка хеширования пароля: %v", err)
		return model.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	// Создание пользователя
	user := model.User{
		ID:        uuid.New().String(),
		Login:     strings.ToLower(strings.TrimSpace(login)), // нормализация
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		log.Printf("AuthService: Ошибка создания пользователя в репозитории: %v", err)
		return model.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	log.Printf("AuthService: Пользователь %s успешно зарегистрирован с ID %s", login, createdUser.ID)
	return createdUser, nil
}
