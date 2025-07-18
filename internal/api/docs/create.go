package docs

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"encoding/json"

	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
	fileserverV1 "github.com/NarthurN/FileServerService/pkg/generated/api/fileserver/v1"
)

// CreateDocument - создание документа
func (a *api) CreateDocument(ctx context.Context, req *fileserverV1.CreateDocumentRequestMultipart) (fileserverV1.CreateDocumentRes, error) {
	// 1. ВАЖНО: Сначала валидируем токен из meta
	token := req.Meta.Token
	if token == "" {
		return &fileserverV1.UnauthorizedError{
			Error: fileserverV1.UnauthorizedErrorError{
				Code: 401,
				Text: "token is required",
			},
		}, nil
	}

	// Валидация токена и получение UserID
	userID, err := a.validateToken(token)
	if err != nil {
		return &fileserverV1.UnauthorizedError{
			Error: fileserverV1.UnauthorizedErrorError{
				Code: 401,
				Text: fmt.Sprintf("invalid token: %v", err),
			},
		}, nil
	}

	name := req.Meta.Name
	if name == "" {
		return &fileserverV1.BadRequestError{
			Error: fileserverV1.BadRequestErrorError{
				Code: 400,
				Text: "name is required",
			},
		}, nil
	}

	meta := req.Meta
	docID := uuid.New()
	createdAt := time.Now().UTC()

	// Куда сохраняем файл (e.g., ./bin/storage/<uuid>)
	storageDir := filepath.Join("bin", "storage")
	if err := os.MkdirAll(storageDir, 0o750); err != nil {
		return &fileserverV1.InternalServerError{
			Error: fileserverV1.InternalServerErrorError{
				Code: 500,
				Text: fmt.Sprintf("create storage dir: %v", err),
			},
		}, nil
	}
	dstPath := filepath.Join(storageDir, docID.String())

	// 2. ИСПРАВЛЕНИЕ: Копируем файл только если он есть
	var actualFilePath string
	if fileData, ok := req.File.Get(); ok {
		dstFile, err := os.Create(dstPath)
		if err != nil {
			return &fileserverV1.InternalServerError{
				Error: fileserverV1.InternalServerErrorError{
					Code: 500,
					Text: fmt.Sprintf("create dst file: %v", err),
				},
			}, nil
		}
		defer dstFile.Close()

		// Копируем содержимое файла
		if _, err := io.Copy(dstFile, fileData.File); err != nil {
			return &fileserverV1.InternalServerError{
				Error: fileserverV1.InternalServerErrorError{
					Code: 500,
					Text: fmt.Sprintf("copy file content: %v", err),
				},
			}, nil
		}
		actualFilePath = dstPath
	}

	// Build business model
	doc := buisnesModel.Document{
		ID:        docID.String(),
		UserID:    userID, // 3. ИСПРАВЛЕНИЕ: Теперь используем реальный UserID
		Name:      name,
		MimeType:  meta.Mime,
		FilePath:  actualFilePath, // 4. ИСПРАВЛЕНИЕ: Путь только если файл есть
		IsFile:    meta.File,
		IsPublic:  meta.Public,
		JSONData:  nil, // will be set below if json provided
		Grants:    meta.Grant,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	// Декодируем опциональную часть JSON в бизнес-модель, если она предоставлена
	if jsonVal, ok := req.JSON.Get(); ok {
		bizJSON := make(map[string]any)
		for k, raw := range jsonVal {
			var v any
			if err := json.Unmarshal(raw, &v); err != nil {
				return &fileserverV1.BadRequestError{
					Error: fileserverV1.BadRequestErrorError{
						Code: 400,
						Text: fmt.Sprintf("invalid json: %v", err),
					},
				}, nil
			}
			bizJSON[k] = v
		}
		doc.JSONData = bizJSON
	}

	// 5. ДОБАВЛЕНИЕ: Валидация - документ должен иметь либо файл, либо JSON
	if !meta.File && doc.JSONData == nil {
		return &fileserverV1.BadRequestError{
			Error: fileserverV1.BadRequestErrorError{
				Code: 400,
				Text: "document must have either file or json data",
			},
		}, nil
	}

	// Сохраняем бизнес-сущность
	if _, err := a.service.CreateDocument(ctx, doc); err != nil {
		return &fileserverV1.InternalServerError{
			Error: fileserverV1.InternalServerErrorError{
				Code: 500,
				Text: fmt.Sprintf("create document: %v", err),
			},
		}, nil
	}

	// Создаем ответ DTO
	responseData := fileserverV1.CreateDocumentResponseData{
		File: name,
	}
	if jsonVal, ok := req.JSON.Get(); ok {
		responseData.JSON = fileserverV1.NewOptCreateDocumentResponseDataJSON(fileserverV1.CreateDocumentResponseDataJSON(jsonVal))
	}

	resp := &fileserverV1.CreateDocumentResponse{
		Data: responseData,
	}

	return resp, nil
}

// validateToken - валидация токена (нужно реализовать)
func (a *api) validateToken(token string) (string, error) {
	// TODO: Реализовать валидацию токена
	// Это может быть:
	// 1. Проверка в базе данных токенов
	// 2. Декодирование JWT токена
	// 3. Проверка в Redis/кеше

	// Пример заглушки:
	if token == "valid-token-123" {
		return "user-123", nil
	}

	return "", fmt.Errorf("invalid token")
}
