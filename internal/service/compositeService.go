package service

import (
	"context"

	"github.com/NarthurN/FileServerService/internal/cache"
	"github.com/NarthurN/FileServerService/internal/config"
	"github.com/NarthurN/FileServerService/internal/model"
	"github.com/NarthurN/FileServerService/internal/repository"
	"github.com/NarthurN/FileServerService/internal/service/auth"
	"github.com/NarthurN/FileServerService/internal/service/docs"
)

// FileServerService - интерфейс сервиса для работы с документами
type DocsService interface {
	// Документы
	CreateDocument(ctx context.Context, doc model.Document) (model.Document, error)
	GetDocument(ctx context.Context, id string) (model.Document, error)
	GetListDocuments(ctx context.Context, userID string) ([]model.Document, error)
	DeleteDocument(ctx context.Context, id string) error

	// Получение документов для пользователя
	GetDocumentsForUser(ctx context.Context, requestUserID, targetUserID string) ([]model.Document, error)
	// Проверка прав доступа к документу
	HasAccessToDocument(ctx context.Context, userID, documentID string) (bool, error)
}

// AuthService - интерфейс сервиса авторизации
type AuthService interface {
	// Регистрация и аутентификация
	RegisterUser(ctx context.Context, adminToken, login, password string) (model.User, error)
	AuthenticateUser(ctx context.Context, login, password string) (string, error)
	ValidateToken(ctx context.Context, tokenValue string) (model.User, error)
	LogoutUser(ctx context.Context, tokenValue string) error

	// Управление токенами
	RefreshToken(ctx context.Context, oldToken string) (string, error)
	GetUserByToken(ctx context.Context, tokenValue string) (model.User, error)

	// Получение пользователя по логину
	GetUserByLogin(ctx context.Context, login string) (model.User, error)
}

type compositeService struct {
	authService AuthService
	docsService DocsService
}

func NewCompositeService(repo repository.FileServerRepository, cfg *config.Config, cacheManager *cache.CacheManager) FileServerService {
	return &compositeService{
		authService: auth.NewService(repo, cfg),
		docsService: docs.NewService(repo, cacheManager),
	}
}

// Методы для работы с документами (делегируем в docsService)
func (s *compositeService) CreateDocument(ctx context.Context, doc model.Document) (model.Document, error) {
	return s.docsService.CreateDocument(ctx, doc)
}

func (s *compositeService) GetDocument(ctx context.Context, id string) (model.Document, error) {
	return s.docsService.GetDocument(ctx, id)
}

func (s *compositeService) GetListDocuments(ctx context.Context, userID string) ([]model.Document, error) {
	return s.docsService.GetListDocuments(ctx, userID)
}

func (s *compositeService) DeleteDocument(ctx context.Context, id string) error {
	return s.docsService.DeleteDocument(ctx, id)
}

func (s *compositeService) GetDocumentsForUser(ctx context.Context, requestUserID, targetUserID string) ([]model.Document, error) {
	return s.docsService.GetDocumentsForUser(ctx, requestUserID, targetUserID)
}

func (s *compositeService) HasAccessToDocument(ctx context.Context, userID, documentID string) (bool, error) {
	return s.docsService.HasAccessToDocument(ctx, userID, documentID)
}

// Методы для работы с аутентификацией (делегируем в authService)
func (s *compositeService) RegisterUser(ctx context.Context, adminToken, login, password string) (model.User, error) {
	return s.authService.RegisterUser(ctx, adminToken, login, password)
}

func (s *compositeService) AuthenticateUser(ctx context.Context, login, password string) (string, error) {
	return s.authService.AuthenticateUser(ctx, login, password)
}

func (s *compositeService) ValidateToken(ctx context.Context, tokenValue string) (model.User, error) {
	return s.authService.ValidateToken(ctx, tokenValue)
}

func (s *compositeService) LogoutUser(ctx context.Context, tokenValue string) error {
	return s.authService.LogoutUser(ctx, tokenValue)
}

func (s *compositeService) RefreshToken(ctx context.Context, oldToken string) (string, error) {
	return s.authService.RefreshToken(ctx, oldToken)
}

func (s *compositeService) GetUserByToken(ctx context.Context, tokenValue string) (model.User, error) {
	return s.authService.GetUserByToken(ctx, tokenValue)
}

func (s *compositeService) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	return s.authService.GetUserByLogin(ctx, login)
}
