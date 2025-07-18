package token

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/NarthurN/FileServerService/internal/model"
	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
	"github.com/jackc/pgx/v5"
)

// GetTokenByValue - получение токена по значению
func (r *Repository) GetTokenByValue(ctx context.Context, tokenValue string) (model.Token, error) {
	query, args, err := r.sb.Select("id", "user_id", "token", "expires_at", "created_at", "is_active").
		From("tokens").
		Where(squirrel.And{
			squirrel.Eq{"token": tokenValue},
			squirrel.Eq{"is_active": true},
			squirrel.Gt{"expires_at": "NOW()"},
		}).
		ToSql()
	if err != nil {
		return model.Token{}, err
	}

	var token model.Token
	if err := r.pool.QueryRow(ctx, query, args...).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.IsActive,
	); err != nil {
		if err == pgx.ErrNoRows {
			return model.Token{}, buisnesModel.ErrNotFound
		}
		return model.Token{}, err
	}

	return token, nil
}

// GetTokensByUserID - получение всех активных токенов пользователя
func (r *Repository) GetTokensByUserID(ctx context.Context, userID string) ([]model.Token, error) {
	query, args, err := r.sb.Select("id", "user_id", "token", "expires_at", "created_at", "is_active").
		From("tokens").
		Where(squirrel.And{
			squirrel.Eq{"user_id": userID},
			squirrel.Eq{"is_active": true},
			squirrel.Gt{"expires_at": "NOW()"},
		}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []model.Token
	for rows.Next() {
		var token model.Token
		if err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Token,
			&token.ExpiresAt,
			&token.CreatedAt,
			&token.IsActive,
		); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}
