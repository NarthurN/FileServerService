package v1

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	fileserverV1 "github.com/NarthurN/FileServerService/pkg/generated/api/fileserver/v1"
	"github.com/go-faster/jx"
)

// GetDocument - получение одного документа
func (a *api) GetDocument(ctx context.Context, params fileserverV1.GetDocumentParams) (fileserverV1.GetDocumentRes, error) {
	log.Printf("🔄 API: Получение документа %s", params.ID)

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

	// Получаем документ
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

	// Проверяем права доступа
	hasAccess, err := a.service.HasAccessToDocument(ctx, user.ID, params.ID)
	if err != nil {
		return &fileserverV1.InternalServerError{
			Error: fileserverV1.InternalServerErrorError{
				Code: 500,
				Text: "🚨 Не удалось проверить права доступа",
			},
		}, nil
	}

	if !hasAccess {
		return &fileserverV1.ForbiddenError{
			Error: fileserverV1.ForbiddenErrorError{
				Code: 403,
				Text: "🚨 Доступ запрещен",
			},
		}, nil
	}

	// Если это файл - возвращаем файл
	if doc.IsFile && doc.FilePath != "" {
		file, err := os.Open(doc.FilePath)
		if err != nil {
			return &fileserverV1.InternalServerError{
				Error: fileserverV1.InternalServerErrorError{
					Code: 500,
					Text: "🚨 Не удалось открыть файл",
				},
			}, nil
		}

		log.Printf("🎉 API: Файл %s успешно получен", params.ID)
		return &fileserverV1.GetDocumentOKApplicationOctetStream{
			Data: file,
		}, nil
	}

	// Если это JSON данные - возвращаем JSON
	if doc.JSONData != nil {
		jsonData := make(map[string]jx.Raw)
		for k, v := range doc.JSONData {
			rawData, _ := json.Marshal(v)
			jsonData[k] = jx.Raw(rawData)
		}

		log.Printf("🎉 API: Документ %s успешно получен", params.ID)
		return &fileserverV1.GetDocumentResponse{
			Data: jsonData,
		}, nil
	}

	log.Printf("🚨 API: Документ %s не найден", params.ID)
	return &fileserverV1.NotFoundError{
		Error: fileserverV1.NotFoundErrorError{
			Code: 404,
			Text: "🚨 Документ не содержит данных",
		},
	}, nil
}

// GetDocumentHead - HEAD запрос для документа
func (a *api) GetDocumentHead(ctx context.Context, params fileserverV1.GetDocumentHeadParams) (fileserverV1.GetDocumentHeadRes, error) {
	log.Printf("🔄 API: HEAD запрос для документа %s", params.ID)
	// Валидация токена
	user, err := a.validateToken(ctx, params.Token)
	if err != nil {
		return &fileserverV1.GetDocumentHeadUnauthorized{}, nil
	}

	// Проверяем существование и права доступа
	hasAccess, err := a.service.HasAccessToDocument(ctx, user.ID, params.ID)
	if err != nil {
		return &fileserverV1.GetDocumentHeadInternalServerError{}, nil
	}

	if !hasAccess {
		return &fileserverV1.GetDocumentHeadForbidden{}, nil
	}

	log.Printf("🎉 API: HEAD запрос для документа %s успешно выполнен", params.ID)
	// HEAD запрос не должен возвращать данные
	return &fileserverV1.GetDocumentHeadOK{}, nil
}
