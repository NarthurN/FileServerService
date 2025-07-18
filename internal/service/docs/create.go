package docs

import (
	"context"
	"fmt"
	"log"
	"strings"

	buisnesModel "github.com/NarthurN/FileServerService/internal/model"
)

// CreateDocument - создание документа с бизнес-логикой
func (s *service) CreateDocument(ctx context.Context, doc buisnesModel.Document) (buisnesModel.Document, error) {
	log.Printf("ServiceLayer: Начало создания документа %s для пользователя %s", doc.Name, doc.UserID)

	// Бизнес-валидация
	if err := s.validateDocumentForCreation(doc); err != nil {
		log.Printf("ServiceLayer: Ошибка валидации документа %s: %v", doc.Name, err)
		return buisnesModel.Document{}, fmt.Errorf("validation failed: %w", err)
	}

	// Проверяем, что пользователь существует
	if _, err := s.repo.GetUserByID(ctx, doc.UserID); err != nil {
		log.Printf("ServiceLayer: Пользователь %s не найден: %v", doc.UserID, err)
		return buisnesModel.Document{}, fmt.Errorf("user not found: %w", err)
	}

	// Нормализация данных
	doc = s.normalizeDocument(doc)

	// Проверяем уникальность имени документа для пользователя
	existingDocs, err := s.repo.GetListDocuments(ctx, doc.UserID)
	if err != nil {
		log.Printf("ServiceLayer: Ошибка получения списка документов: %v", err)
		return buisnesModel.Document{}, fmt.Errorf("failed to check existing documents: %w", err)
	}

	for _, existing := range existingDocs {
		if strings.EqualFold(existing.Name, doc.Name) {
			log.Printf("ServiceLayer: Документ с именем %s уже существует", doc.Name)
			return buisnesModel.Document{}, fmt.Errorf("document with name '%s' already exists", doc.Name)
		}
	}

	// Валидируем grants (проверяем, что пользователи существуют)
	if err := s.validateGrants(ctx, doc.Grants); err != nil {
		log.Printf("ServiceLayer: Ошибка валидации grants: %v", err)
		return buisnesModel.Document{}, fmt.Errorf("invalid grants: %w", err)
	}

	// Создаем документ
	createdDoc, err := s.repo.CreateDocument(ctx, doc)
	if err != nil {
		log.Printf("ServiceLayer: Ошибка создания документа в репозитории: %v", err)
		return buisnesModel.Document{}, fmt.Errorf("failed to create document: %w", err)
	}

	log.Printf("ServiceLayer: Документ %s успешно создан с ID %s", createdDoc.Name, createdDoc.ID)
	return createdDoc, nil
}
