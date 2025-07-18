package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/NarthurN/FileServerService/internal/config"
	"github.com/NarthurN/FileServerService/internal/model"
	"github.com/NarthurN/FileServerService/internal/repository"
)

type Service struct {
	repo   repository.FileServerRepository
	config *config.Config
}

func NewService(repo repository.FileServerRepository, cfg *config.Config) *Service {
	return &Service{
		repo:   repo,
		config: cfg,
	}
}

// Вспомогательные методы с бизнес-логикой

func (s *Service) validateAdminToken(token string) error {
	if token != s.config.Auth.AdminToken {
		return fmt.Errorf("admin token mismatch")
	}
	return nil
}

func (s *Service) validateLogin(login string) error {
	login = strings.TrimSpace(login)

	if len(login) < 8 {
		return fmt.Errorf("login must be at least 8 characters long")
	}

	if len(login) > 50 {
		return fmt.Errorf("login must be no more than 50 characters long")
	}

	// Только латиница и цифры
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, login)
	if !matched {
		return fmt.Errorf("login must contain only latin letters and digits")
	}

	return nil
}

func (s *Service) validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return fmt.Errorf("password must be no more than 128 characters long")
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasDigit   = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

func (s *Service) validateAuthInput(login, password string) error {
	if strings.TrimSpace(login) == "" {
		return fmt.Errorf("login is required")
	}
	if strings.TrimSpace(password) == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

func (s *Service) checkLoginUniqueness(ctx context.Context, login string) error {
	normalizedLogin := strings.ToLower(strings.TrimSpace(login))
	_, err := s.repo.GetUserByLogin(ctx, normalizedLogin)
	if err == nil {
		return fmt.Errorf("user with login '%s' already exists", login)
	}
	if err != model.ErrNotFound {
		return fmt.Errorf("failed to check login uniqueness: %w", err)
	}
	return nil
}

func (s *Service) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (s *Service) verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *Service) generateSecureToken() (string, error) {
	// Генерируем 32 случайных байта
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Добавляем временную метку для уникальности
	timestamp := time.Now().UnixNano()
	timestampBytes := []byte(fmt.Sprintf("%d", timestamp))

	// Объединяем и хешируем
	combined := append(bytes, timestampBytes...)
	hash := sha256.Sum256(combined)

	return hex.EncodeToString(hash[:]), nil
}

func (s *Service) getTokenLifetime() time.Duration {
	// Настраиваемое время жизни токена (по умолчанию 24 часа)
	if s.config.Auth.TokenLifetime != 0 {
		return s.config.Auth.TokenLifetime
	}
	return 24 * time.Hour
}

func (s *Service) cleanupOldTokens(ctx context.Context, userID string) error {
	// Можно настроить: деактивировать все старые токены или оставить несколько активных
	return s.repo.DeactivateUserTokens(ctx, userID)
}
