package user

import (
	"context"
	"log"

	"github.com/NarthurN/FileServerService/internal/model"
)

// CreateUser - создание нового пользователя
func (r *Repository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	log.Printf("RepLayer: Начало создания пользователя %s\n", user.Login)

	query, args, err := r.sb.Insert("users").
		Columns("id", "login", "password_hash", "created_at", "updated_at").
		Values(user.ID, user.Login, user.Password, user.CreatedAt, user.UpdatedAt).
		ToSql()
	if err != nil {
		log.Printf("RepLayer: ошибка подготовки запроса создания пользователя %s: %v\n", user.Login, err)
		return model.User{}, err
	}

	log.Printf("RepLayer: Запрос для создания пользователя %s подготовлен\n", user.Login)

	if _, err := r.pool.Exec(ctx, query, args...); err != nil {
		log.Printf("RepLayer: ошибка создания пользователя %s: %v\n", user.Login, err)
		return model.User{}, err
	}

	log.Printf("RepLayer: Пользователь %s создан\n", user.Login)
	return user, nil
}
