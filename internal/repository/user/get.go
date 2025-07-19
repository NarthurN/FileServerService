package user

import (
	"context"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/NarthurN/FileServerService/internal/model"
	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
	"github.com/jackc/pgx/v5"
)

// GetUserByLogin - получение пользователя по логину
func (r *Repository) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	query, args, err := r.sb.Select("id", "login", "password_hash", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"login": login}).
		ToSql()
	if err != nil {
		return model.User{}, err
	}

	var user model.User
	if err := r.pool.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return model.User{}, buisnesModel.ErrNotFound
		}
		return model.User{}, err
	}

	return user, nil
}

// GetUserByID - получение пользователя по ID
func (r *Repository) GetUserByID(ctx context.Context, userID string) (model.User, error) {
	log.Printf("Repository: Получение пользователя по ID %s", userID)

	query, args, err := r.sb.Select("id", "login", "password_hash", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": userID}).
		ToSql()
	if err != nil {
		log.Printf("Repository: Ошибка создания SQL запроса: %v", err)
		return model.User{}, err
	}

	var user model.User
	if err := r.pool.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("Repository: Пользователь %s не найден", userID)
			return model.User{}, buisnesModel.ErrNotFound
		}
		log.Printf("Repository: Ошибка сканирования пользователя: %v", err)
		return model.User{}, err
	}

	log.Printf("Repository: Пользователь %s найден: %s", userID, user.Login)
	return user, nil
}
