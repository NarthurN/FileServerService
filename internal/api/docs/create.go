package docs

import (
	"context"
	"fmt"
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
	name := req.Meta.Name
	if name == "" {
		return &fileserverV1.BadRequestError{
			Code:    400,
			Message: "name is required",
		}, nil
	}

	meta := req.Meta
	docID := uuid.New()
	createdAt := time.Now().UTC()

	// Куда сохраняем файл (e.g., ./bin/storage/<uuid>)
	storageDir := filepath.Join("bin", "storage")
	if err := os.MkdirAll(storageDir, 0o750); err != nil {
		return &fileserverV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("create storage dir: %v", err),
		}, nil
	}
	dstPath := filepath.Join(storageDir, docID.String())

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return &fileserverV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("create dst file: %v", err),
		}, nil
	}
	defer dstFile.Close()

	// Build business model
	doc := buisnesModel.Document{
		ID:        docID.String(),
		UserID:    "", // TODO: extract from auth context
		Name:      name,
		MimeType:  meta.Mime,
		FilePath:  dstPath,
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
					Code:    400,
					Message: fmt.Sprintf("invalid json: %v", err),
				}, nil
			}
			bizJSON[k] = v
		}
		doc.JSONData = bizJSON
	}

	// Сохраняем бизнес-сущность
	if _, err := a.docsService.CreateDocument(ctx, doc); err != nil {
		return &fileserverV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("create document: %v", err),
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
