package token

import (
	"context"
	"log"

	"github.com/NarthurN/FileServerService/internal/model"
)

// CreateToken - создание токена
func (r *Repository) CreateToken(ctx context.Context, token model.Token) (model.Token, error) {
	log.Printf("RepLayer: Начало создания токена для пользователя %s\n", token.UserID)

	query, args, err := r.sb.Insert("tokens").
		Columns("id", "user_id", "token", "expires_at", "created_at", "is_active").
		Values(token.ID, token.UserID, token.Token, token.ExpiresAt, token.CreatedAt, token.IsActive).
		ToSql()
	if err != nil {
		log.Printf("RepLayer: ошибка подготовки запроса создания токена: %v\n", err)
		return model.Token{}, err
	}

	log.Printf("RepLayer: Запрос для создания токена подготовлен\n")

	if _, err := r.pool.Exec(ctx, query, args...); err != nil {
		log.Printf("RepLayer: ошибка создания токена: %v\n", err)
		return model.Token{}, err
	}

	log.Printf("RepLayer: Токен создан для пользователя %s\n", token.UserID)
	return token, nil
}
