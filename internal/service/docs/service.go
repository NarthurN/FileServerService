package docs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/NarthurN/FileServerService/internal/cache"
	"github.com/NarthurN/FileServerService/internal/model"
	"github.com/NarthurN/FileServerService/internal/repository"
)

type service struct {
	repo         repository.FileServerRepository
	cacheManager *cache.CacheManager
}

func NewService(repo repository.FileServerRepository, cacheManager *cache.CacheManager) *service {
	return &service{
		repo:         repo,
		cacheManager: cacheManager,
	}
}

// Вспомогательные методы с бизнес-логикой
func (s *service) validateDocumentForCreation(doc model.Document) error {
	if doc.Name == "" {
		return fmt.Errorf("document name is required")
	}

	if len(doc.Name) > 255 {
		return fmt.Errorf("document name too long (max 255 characters)")
	}

	if doc.UserID == "" {
		return fmt.Errorf("user ID is required")
	}

	if doc.MimeType == "" {
		return fmt.Errorf("MIME type is required")
	}

	// Если это файл, должен быть указан путь к файлу
	if doc.IsFile && doc.FilePath == "" {
		return fmt.Errorf("file path is required for file documents")
	}

	// Если это не файл, должны быть JSON данные
	if !doc.IsFile && doc.JSONData == nil {
		return fmt.Errorf("JSON data is required for non-file documents")
	}

	return nil
}

func (s *service) normalizeDocument(doc model.Document) model.Document {
	// Нормализуем имя файла
	doc.Name = strings.TrimSpace(doc.Name)

	// Нормализуем MIME type
	doc.MimeType = strings.ToLower(strings.TrimSpace(doc.MimeType))

	// Нормализуем grants (убираем дубликаты и пустые значения)
	normalizedGrants := make([]string, 0, len(doc.Grants))
	seenGrants := make(map[string]bool)

	for _, grant := range doc.Grants {
		grant = strings.TrimSpace(grant)
		if grant != "" && !seenGrants[grant] {
			normalizedGrants = append(normalizedGrants, grant)
			seenGrants[grant] = true
		}
	}
	doc.Grants = normalizedGrants

	return doc
}

func (s *service) validateGrants(ctx context.Context, grants []string) error {
	for _, login := range grants {
		if login == "" {
			continue // пропускаем пустые
		}

		// Проверяем, что пользователь с таким логином существует
		if _, err := s.repo.GetUserByLogin(ctx, login); err != nil {
			if err == model.ErrNotFound {
				return fmt.Errorf("user with login '%s' not found", login)
			}
			return fmt.Errorf("failed to validate user '%s': %w", login, err)
		}
	}
	return nil
}

func (s *service) hasAccessToDocument(doc model.Document, userLogin string) bool {
	// Публичные документы доступны всем
	if doc.IsPublic {
		return true
	}

	// Проверяем grants
	for _, grantedLogin := range doc.Grants {
		if grantedLogin == userLogin {
			return true
		}
	}

	return false
}

func (s *service) sortDocuments(docs []model.Document) []model.Document {
	// Согласно заданию - сортировка по имени и дате создания
	// Реализуем простую сортировку
	log.Printf("ServiceLayer: Сортировка %d документов", len(docs))

	if len(docs) <= 1 {
		return docs
	}

	// Создаем копию для сортировки
	sortedDocs := make([]model.Document, len(docs))
	copy(sortedDocs, docs)

	// Простая bubble sort (для production лучше использовать sort.Slice)
	for i := 0; i < len(sortedDocs)-1; i++ {
		for j := 0; j < len(sortedDocs)-i-1; j++ {
			if s.shouldSwapDocuments(sortedDocs[j], sortedDocs[j+1]) {
				sortedDocs[j], sortedDocs[j+1] = sortedDocs[j+1], sortedDocs[j]
			}
		}
	}

	return sortedDocs
}

func (s *service) shouldSwapDocuments(doc1, doc2 model.Document) bool {
	// Сначала сортируем по имени
	if doc1.Name != doc2.Name {
		return doc1.Name > doc2.Name
	}

	// Если имена одинаковые, сортируем по дате создания (новые первыми)
	return doc1.CreatedAt.Before(doc2.CreatedAt)
}
