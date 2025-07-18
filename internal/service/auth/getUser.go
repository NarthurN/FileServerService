package auth

import (
	"context"

	"github.com/NarthurN/FileServerService/internal/model"
)

// GetUserByToken - получение пользователя по токену
func (s *Service) GetUserByToken(ctx context.Context, tokenValue string) (model.User, error) {
	return s.ValidateToken(ctx, tokenValue)
}
