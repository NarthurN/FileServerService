package service

import (
	"context"

	"github.com/NarthurN/FileServerService/internal/model"
)

// FileServerService - интерфейс сервиса для работы с документами
type FileServerService interface {
	// Документы
	CreateDocument(ctx context.Context, doc model.Document) (model.Document, error)
	GetDocument(ctx context.Context, id string) (model.Document, error)
	GetListDocuments(ctx context.Context, userID string) ([]model.Document, error)
	DeleteDocument(ctx context.Context, id string) error

	// Получение документов для пользователя
	GetDocumentsForUser(ctx context.Context, requestUserID, targetUserID string) ([]model.Document, error)

	// Регистрация и аутентификация
	RegisterUser(ctx context.Context, adminToken, login, password string) (model.User, error)
	AuthenticateUser(ctx context.Context, login, password string) (string, error)
	ValidateToken(ctx context.Context, tokenValue string) (model.User, error)
	LogoutUser(ctx context.Context, tokenValue string) error

	// Управление токенами
	RefreshToken(ctx context.Context, oldToken string) (string, error)
	GetUserByToken(ctx context.Context, tokenValue string) (model.User, error)
}
