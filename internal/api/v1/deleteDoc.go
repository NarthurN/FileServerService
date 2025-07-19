package v1

import (
	"context"
	"log"
	"os"
	"strings"

	fileserverV1 "github.com/NarthurN/FileServerService/pkg/generated/api/fileserver/v1"
)

// DeleteDocument - удаление документа
func (a *api) DeleteDocument(ctx context.Context, params fileserverV1.DeleteDocumentParams) (fileserverV1.DeleteDocumentRes, error) {
	log.Printf("🔄 API: Удаление документа %s", params.ID)

	// Валидация токена
	user, err := a.validateToken(ctx, params.Token)
	if err != nil {
		return &fileserverV1.UnauthorizedError{
			Error: fileserverV1.UnauthorizedErrorError{
				Code: 401,
				Text: "🚨 Неверный токен",
			},
		}, nil
	}

	// Получаем документ для проверки прав
	doc, err := a.service.GetDocument(ctx, params.ID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return &fileserverV1.NotFoundError{
				Error: fileserverV1.NotFoundErrorError{
					Code: 404,
					Text: "🚨 Документ не найден",
				},
			}, nil
		}

		return &fileserverV1.InternalServerError{
			Error: fileserverV1.InternalServerErrorError{
				Code: 500,
				Text: "🚨 Не удалось получить документ",
			},
		}, nil
	}

	// Проверяем, что пользователь является владельцем
	if doc.UserID != user.ID {
		return &fileserverV1.ForbiddenError{
			Error: fileserverV1.ForbiddenErrorError{
				Code: 403,
				Text: "🚨 Только владелец документа может удалить его",
			},
		}, nil
	}

	// Удаляем физический файл, если он есть
	if doc.IsFile && doc.FilePath != "" {
		if err := os.Remove(doc.FilePath); err != nil {
			log.Printf("🚨 API: Предупреждение - не удалось удалить файл %s: %v", doc.FilePath, err)
		}
	}

	// Удаляем документ через сервис
	if err := a.service.DeleteDocument(ctx, params.ID); err != nil {
		return &fileserverV1.InternalServerError{
			Error: fileserverV1.InternalServerErrorError{
				Code: 500,
				Text: "🚨 Не удалось удалить документ",
			},
		}, nil
	}

	log.Printf("🎉 API: Документ %s успешно удален", params.ID)

	// Формируем ответ согласно заданию
	response := make(fileserverV1.DeleteDocumentResponseResponse)
	response[params.ID] = true

	return &fileserverV1.DeleteDocumentResponse{
		Response: response,
	}, nil
}
