package v1

import (
	"context"
	"fmt"
	"log"

	"github.com/NarthurN/FileServerService/internal/model"
	"github.com/NarthurN/FileServerService/internal/service"
	fileserverV1 "github.com/NarthurN/FileServerService/pkg/generated/api/fileserver/v1"
)

var _ fileserverV1.Handler = (*api)(nil)

type api struct {
	//fileserverV1.UnimplementedHandler

	service service.FileServerService
}

func NewAPI(service service.FileServerService) *api {
	return &api{
		service: service,
	}
}

// validateToken - валидация токена и получение пользователя
func (a *api) validateToken(ctx context.Context, token string) (model.User, error) {
	if token == "" {
		return model.User{}, fmt.Errorf("token is required")
	}

	user, err := a.service.ValidateToken(ctx, token)
	if err != nil {
		return model.User{}, fmt.Errorf("invalid token: %w", err)
	}

	return user, nil
}

// getUserByLogin - получение пользователя по логину
func (a *api) getUserByLogin(ctx context.Context, login string) (model.User, error) {
	log.Printf("API: Поиск пользователя по логину: %s", login)
	// Используем сервис для получения пользователя по логину
	user, err := a.service.GetUserByLogin(ctx, login)
	if err != nil {
		log.Printf("API: Пользователь %s не найден: %v", login, err)
		return model.User{}, fmt.Errorf("user not found: %w", err)
	}
	log.Printf("API: Пользователь %s найден, ID: %s", login, user.ID)
	return user, nil
}

// filterDocuments - фильтрация документов по ключу и значению
func (a *api) filterDocuments(docs []model.Document, key, value string) []model.Document {
	var filtered []model.Document

	for _, doc := range docs {
		match := false
		switch key {
		case "name":
			match = doc.Name == value
		case "mime":
			match = doc.MimeType == value
		case "public":
			match = (value == "true" && doc.IsPublic) || (value == "false" && !doc.IsPublic)
		case "file":
			match = (value == "true" && doc.IsFile) || (value == "false" && !doc.IsFile)
		case "created":
			// Простое сравнение даты в формате строки
			match = doc.CreatedAt.Format("2006-01-02 15:04:05") == value
		default:
			// Если ключ не распознан, возвращаем все документы
			return docs
		}

		if match {
			filtered = append(filtered, doc)
		}
	}

	return filtered
}
