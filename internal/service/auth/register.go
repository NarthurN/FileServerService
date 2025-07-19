package auth

import (
	"context"
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
		return model.User{}, model.NewAuthError("Неверный админский токен", model.ErrInvalidAdminToken)
	}

	// Валидация логина согласно заданию
	if err := s.validateLogin(login); err != nil {
		log.Printf("AuthService: Неверный логин %s: %v", login, err)
		return model.User{}, model.NewValidationError("Неверный логин", err)
	}

	// Валидация пароля согласно заданию
	if err := s.validatePassword(password); err != nil {
		log.Printf("AuthService: Неверный пароль для пользователя %s: %v", login, err)
		return model.User{}, model.NewValidationError("Неверный пароль", err)
	}

	// Бизнес-правило: проверяем уникальность логина
	if err := s.checkLoginUniqueness(ctx, login); err != nil {
		log.Printf("AuthService: Логин %s уже занят", login)
		return model.User{}, model.NewValidationError("Логин уже занят", model.ErrLoginAlreadyExists)
	}

	// Бизнес-логика: хеширование пароля
	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		log.Printf("AuthService: Ошибка хеширования пароля: %v", err)
		return model.User{}, model.NewBusinessError("Ошибка хеширования пароля", err)
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
		return model.User{}, model.NewBusinessError("Ошибка создания пользователя в репозитории", err)
	}

	log.Printf("AuthService: Пользователь %s успешно зарегистрирован с ID %s", login, createdUser.ID)
	return createdUser, nil
}
