package v1

import (
	"context"
	"log"

	"github.com/NarthurN/FileServerService/internal/model"
	fileserverV1 "github.com/NarthurN/FileServerService/pkg/generated/api/fileserver/v1"
)

// ListDocuments - получение списка документов
func (a *api) ListDocuments(ctx context.Context, params fileserverV1.ListDocumentsParams) (fileserverV1.ListDocumentsRes, error) {
	log.Printf("🔄 API: Получение списка документов")

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

	// Определяем, чьи документы получать
	var docs []model.Document
	if loginParam, ok := params.Login.Get(); ok && loginParam != "" {
		// Получаем документы другого пользователя (только те, к которым есть доступ)
		targetUser, err := a.getUserByLogin(ctx, loginParam)
		if err != nil {
			return &fileserverV1.InternalServerError{
				Error: fileserverV1.InternalServerErrorError{
					Code: 500,
					Text: "🚨 Не удалось получить документы",
				},
			}, nil
		}

		docs, err = a.service.GetDocumentsForUser(ctx, user.ID, targetUser.ID)
		if err != nil {
			return &fileserverV1.InternalServerError{
				Error: fileserverV1.InternalServerErrorError{
					Code: 500,
					Text: "🚨 Не удалось получить документы",
				},
			}, nil
		}
	} else {
		// Получаем собственные документы
		docs, err = a.service.GetListDocuments(ctx, user.ID)
		if err != nil {
			return &fileserverV1.InternalServerError{
				Error: fileserverV1.InternalServerErrorError{
					Code: 500,
					Text: "🚨 Не удалось получить документы",
				},
			}, nil
		}
	}

	// Применяем фильтрацию если указана
	if keyParam, keyOk := params.Key.Get(); keyOk {
		if valueParam, valueOk := params.Value.Get(); valueOk && valueParam != "" {
			docs = a.filterDocuments(docs, string(keyParam), valueParam)
		}
	}

	// Применяем лимит если указан
	if limitParam, ok := params.Limit.Get(); ok && limitParam > 0 {
		if len(docs) > limitParam {
			docs = docs[:limitParam]
		}
	}

	// Конвертируем в DTO для ответа
	docDTOs := make([]fileserverV1.DocumentDto, 0, len(docs))
	for _, doc := range docs {
		dto := fileserverV1.DocumentDto{
			ID:      doc.ID,
			Name:    doc.Name,
			Mime:    doc.MimeType,
			File:    doc.IsFile,
			Public:  doc.IsPublic,
			Created: doc.CreatedAt.Format("2006-01-02 15:04:05"),
			Grant:   doc.Grants,
		}
		docDTOs = append(docDTOs, dto)
	}

	log.Printf("🎉 API: Найдено %d документов", len(docDTOs))

	return &fileserverV1.ListDocumentsResponse{
		Data: fileserverV1.ListDocumentsResponseData{
			Docs: docDTOs,
		},
	}, nil
}

// ListDocumentsHead - HEAD запрос для списка документов
func (a *api) ListDocumentsHead(ctx context.Context, params fileserverV1.ListDocumentsHeadParams) (fileserverV1.ListDocumentsHeadRes, error) {
	// Валидация токена
	_, err := a.validateToken(ctx, params.Token)
	if err != nil {
		return &fileserverV1.ListDocumentsHeadUnauthorized{}, nil
	}

	// HEAD запрос не должен возвращать данные согласно заданию
	return &fileserverV1.ListDocumentsHeadOK{}, nil
}
