package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
	"github.com/NarthurN/FileServerService/internal/repository/doc"
	"github.com/NarthurN/FileServerService/internal/repository/token"
	"github.com/NarthurN/FileServerService/internal/repository/user"
)

var _ FileServerRepository = (*CompositeRepository)(nil)

type docRepository interface {
	CreateDocument(ctx context.Context, doc buisnesModel.Document) (buisnesModel.Document, error)
	GetDocument(ctx context.Context, id string) (buisnesModel.Document, error)
	GetListDocuments(ctx context.Context, userID string) ([]buisnesModel.Document, error)
	DeleteDocument(ctx context.Context, id string) error
}

type userRepository interface {
	CreateUser(ctx context.Context, user buisnesModel.User) (buisnesModel.User, error)
	GetUserByLogin(ctx context.Context, login string) (buisnesModel.User, error)
	GetUserByID(ctx context.Context, userID string) (buisnesModel.User, error)
}

type tokenRepository interface {
	CreateToken(ctx context.Context, token buisnesModel.Token) (buisnesModel.Token, error)
	GetTokenByValue(ctx context.Context, tokenValue string) (buisnesModel.Token, error)
	DeactivateToken(ctx context.Context, tokenValue string) error
	DeactivateUserTokens(ctx context.Context, userID string) error
}

// CompositeRepository - композитный репозиторий, объединяющий все репозитории
type CompositeRepository struct {
	userRepo  userRepository
	docRepo   docRepository
	tokenRepo tokenRepository
}

func NewCompositeRepository(pool *pgxpool.Pool) *CompositeRepository {
	return &CompositeRepository{
		userRepo:  user.NewRepository(pool),
		docRepo:   doc.NewRepository(pool),
		tokenRepo: token.NewRepository(pool),
	}
}

// Методы для работы с документами (делегируем в docRepo)
func (r *CompositeRepository) CreateDocument(ctx context.Context, doc buisnesModel.Document) (buisnesModel.Document, error) {
	return r.docRepo.CreateDocument(ctx, doc)
}

func (r *CompositeRepository) GetDocument(ctx context.Context, id string) (buisnesModel.Document, error) {
	return r.docRepo.GetDocument(ctx, id)
}

func (r *CompositeRepository) GetListDocuments(ctx context.Context, userID string) ([]buisnesModel.Document, error) {
	return r.docRepo.GetListDocuments(ctx, userID)
}

func (r *CompositeRepository) DeleteDocument(ctx context.Context, id string) error {
	return r.docRepo.DeleteDocument(ctx, id)
}

// Методы для работы с пользователями (делегируем в userRepo)
func (r *CompositeRepository) CreateUser(ctx context.Context, user buisnesModel.User) (buisnesModel.User, error) {
	return r.userRepo.CreateUser(ctx, user)
}

func (r *CompositeRepository) GetUserByLogin(ctx context.Context, login string) (buisnesModel.User, error) {
	return r.userRepo.GetUserByLogin(ctx, login)
}

func (r *CompositeRepository) GetUserByID(ctx context.Context, userID string) (buisnesModel.User, error) {
	return r.userRepo.GetUserByID(ctx, userID)
}

// Методы для работы с токенами (делегируем в tokenRepo)
func (r *CompositeRepository) CreateToken(ctx context.Context, token buisnesModel.Token) (buisnesModel.Token, error) {
	return r.tokenRepo.CreateToken(ctx, token)
}

func (r *CompositeRepository) GetTokenByValue(ctx context.Context, tokenValue string) (buisnesModel.Token, error) {
	return r.tokenRepo.GetTokenByValue(ctx, tokenValue)
}

func (r *CompositeRepository) DeactivateToken(ctx context.Context, tokenValue string) error {
	return r.tokenRepo.DeactivateToken(ctx, tokenValue)
}

func (r *CompositeRepository) DeactivateUserTokens(ctx context.Context, userID string) error {
	return r.tokenRepo.DeactivateUserTokens(ctx, userID)
}
