package repository

import (
	"context"

	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
)

type FileServerRepository interface {
	// Документы
	CreateDocument(ctx context.Context, doc buisnesModel.Document) (buisnesModel.Document, error)
	GetDocument(ctx context.Context, id string) (buisnesModel.Document, error)
	GetListDocuments(ctx context.Context, userID string) ([]buisnesModel.Document, error)
	DeleteDocument(ctx context.Context, id string) error

	// Пользователи
	CreateUser(ctx context.Context, user buisnesModel.User) (buisnesModel.User, error)
	GetUserByLogin(ctx context.Context, login string) (buisnesModel.User, error)
	GetUserByID(ctx context.Context, userID string) (buisnesModel.User, error)

	// Токены
	CreateToken(ctx context.Context, token buisnesModel.Token) (buisnesModel.Token, error)
	GetTokenByValue(ctx context.Context, tokenValue string) (buisnesModel.Token, error)
	DeactivateToken(ctx context.Context, tokenValue string) error
	DeactivateUserTokens(ctx context.Context, userID string) error
}
