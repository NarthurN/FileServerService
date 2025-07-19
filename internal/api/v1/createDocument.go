package v1

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"encoding/json"

	"github.com/NarthurN/FileServerService/internal/model"
	fileserverV1 "github.com/NarthurN/FileServerService/pkg/generated/api/fileserver/v1"
)

// CreateDocument - создание документа
func (a *api) CreateDocument(ctx context.Context, req *fileserverV1.CreateDocumentRequestMultipart) (fileserverV1.CreateDocumentRes, error) {
	log.Printf("🔄 API: Создание документа %s", req.Meta.Name)

	// Валидация токена
	user, err := a.validateToken(ctx, req.Meta.Token)
	if err != nil {
		log.Printf("🚨 API: Ошибка валидации токена: %v", err)
		return &fileserverV1.UnauthorizedError{
			Error: fileserverV1.UnauthorizedErrorError{
				Code: 401,
				Text: fmt.Sprintf("🚨 Токен %s: %v", req.Meta.Token, err),
			},
		}, nil
	}

	// Валидация имени документа
	if req.Meta.Name == "" {
		return &fileserverV1.BadRequestError{
			Error: fileserverV1.BadRequestErrorError{
				Code: 400,
				Text: "🚨 Имя документа не может быть пустым",
			},
		}, nil
	}

	docID := uuid.New().String()

	// Создание файла
	var filePath string
	if req.Meta.File {
		if fileData, ok := req.File.Get(); ok {
			storageDir := filepath.Join("bin", "storage")
			if err := os.MkdirAll(storageDir, 0755); err != nil {
				log.Printf("🚨 API: Ошибка создания директории: %v", err)
				return &fileserverV1.InternalServerError{
					Error: fileserverV1.InternalServerErrorError{
						Code: 500,
						Text: fmt.Sprintf("🚨 Ошибка создания директории: %v", err),
					},
				}, nil
			}

			filePath = filepath.Join(storageDir, docID)
			dstFile, err := os.Create(filePath)
			if err != nil {
				log.Printf("🚨 API: Ошибка создания файла: %v", err)
				return &fileserverV1.InternalServerError{
					Error: fileserverV1.InternalServerErrorError{
						Code: 500,
						Text: fmt.Sprintf("🚨 Ошибка создания файла: %v", err),
					},
				}, nil
			}
			defer dstFile.Close()

			if _, err := io.Copy(dstFile, fileData.File); err != nil {
				log.Printf("🚨 API: Ошибка копирования файла: %v", err)
				return &fileserverV1.InternalServerError{
					Error: fileserverV1.InternalServerErrorError{
						Code: 500,
						Text: fmt.Sprintf("🚨 Ошибка копирования файла: %v", err),
					},
				}, nil
			}
		} else if req.Meta.File {
			log.Printf("🚨 API: Файл не найден")
			return &fileserverV1.BadRequestError{
				Error: fileserverV1.BadRequestErrorError{
					Code: 400,
					Text: "🚨 Файл обязателен если мета.file = true",
				},
			}, nil
		}
	}

	doc := model.Document{
		ID:        docID,
		UserID:    user.ID,
		Name:      req.Meta.Name,
		MimeType:  req.Meta.Mime,
		FilePath:  filePath,
		IsFile:    req.Meta.File,
		IsPublic:  req.Meta.Public,
		JSONData:  nil,
		Grants:    req.Meta.Grant,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Обрабатываем JSON данные (если есть)
	if jsonVal, ok := req.JSON.Get(); ok {
		bizJSON := make(map[string]any)
		for k, raw := range jsonVal {
			var v any
			if err := json.Unmarshal(raw, &v); err != nil {
				log.Printf("🚨 API: Ошибка парсинга JSON: %v", err)
				return &fileserverV1.BadRequestError{
					Error: fileserverV1.BadRequestErrorError{
						Code: 400,
						Text: fmt.Sprintf("🚨 Ошибка парсинга JSON: %v", err),
					},
				}, nil
			}
			bizJSON[k] = v
		}
		doc.JSONData = bizJSON
	}

	// Сохраняем документ через сервис
	createDoc, err := a.service.CreateDocument(ctx, doc)
	if err != nil {
		log.Printf("🚨 API: Ошибка создания документа: %v", err)
		if filePath != "" {
			os.Remove(filePath)
		}
		return &fileserverV1.BadRequestError{
			Error: fileserverV1.BadRequestErrorError{
				Code: 400,
				Text: fmt.Sprintf("🚨 Ошибка парсинга JSON: %v", err),
			},
		}, nil
	}

	// Формируем ответ
	responseData := fileserverV1.CreateDocumentResponseData{
		File: createDoc.Name,
	}
	if jsonVal, ok := req.JSON.Get(); ok {
		responseData.JSON = fileserverV1.NewOptCreateDocumentResponseDataJSON(
			fileserverV1.CreateDocumentResponseDataJSON(jsonVal),
		)
	}
	log.Printf("🎉 API: Документ %s успешно создан", docID)
	return &fileserverV1.CreateDocumentResponse{
		Data: responseData,
	}, nil
}
