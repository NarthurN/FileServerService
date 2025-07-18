package token

import (
	"context"
	"log"

	"github.com/Masterminds/squirrel"
)

// DeactivateToken - деактивация токена
func (r *Repository) DeactivateToken(ctx context.Context, tokenValue string) error {
	log.Printf("RepLayer: Начало деактивации токена\n")

	query, args, err := r.sb.Update("tokens").
		Set("is_active", false).
		Where(squirrel.Eq{"token": tokenValue}).
		ToSql()
	if err != nil {
		log.Printf("RepLayer: ошибка подготовки запроса деактивации токена: %v\n", err)
		return err
	}

	if _, err := r.pool.Exec(ctx, query, args...); err != nil {
		log.Printf("RepLayer: ошибка деактивации токена: %v\n", err)
		return err
	}

	log.Printf("RepLayer: Токен деактивирован\n")
	return nil
}

// DeactivateUserTokens - деактивация всех токенов пользователя
func (r *Repository) DeactivateUserTokens(ctx context.Context, userID string) error {
	log.Printf("RepLayer: Начало деактивации всех токенов пользователя %s\n", userID)

	query, args, err := r.sb.Update("tokens").
		Set("is_active", false).
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		log.Printf("RepLayer: ошибка подготовки запроса деактивации токенов пользователя: %v\n", err)
		return err
	}

	if _, err := r.pool.Exec(ctx, query, args...); err != nil {
		log.Printf("RepLayer: ошибка деактивации токенов пользователя: %v\n", err)
		return err
	}

	log.Printf("RepLayer: Все токены пользователя %s деактивированы\n", userID)
	return nil
}

// DeactivateExpiredTokens - деактивация истекших токенов (полезно для очистки)
func (r *Repository) DeactivateExpiredTokens(ctx context.Context) error {
	log.Printf("RepLayer: Начало деактивации истекших токенов\n")

	query, args, err := r.sb.Update("tokens").
		Set("is_active", false).
		Where(squirrel.And{
			squirrel.Eq{"is_active": true},
			squirrel.Lt{"expires_at": "NOW()"},
		}).
		ToSql()
	if err != nil {
		log.Printf("RepLayer: ошибка подготовки запроса деактивации истекших токенов: %v\n", err)
		return err
	}

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("RepLayer: ошибка деактивации истекших токенов: %v\n", err)
		return err
	}

	log.Printf("RepLayer: Деактивировано %d истекших токенов\n", result.RowsAffected())
	return nil
}
