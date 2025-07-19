package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

// CacheManager управляет кэшированием данных
type CacheManager struct {
	cache Cache
}

// NewCacheManager создает новый кэш-менеджер
func NewCacheManager(capacity int64) (*CacheManager, error) {
	cache := NewLRUCache(capacity)

	return &CacheManager{
		cache: cache,
	}, nil
}

// CacheKey генерирует ключ кэша для различных типов данных
type CacheKey struct {
	Type       string
	UserID     string
	DocumentID string
	Filter     string
	Limit      int
}

// GenerateKey генерирует уникальный ключ кэша
func (ck *CacheKey) GenerateKey() string {
	data := fmt.Sprintf("%s:%s:%s:%s:%d",
		ck.Type, ck.UserID, ck.DocumentID, ck.Filter, ck.Limit)

	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// DocumentListKey создает ключ для списка документов
func DocumentListKey(userID, filterKey, filterValue string, limit int) string {
	filter := ""
	if filterKey != "" && filterValue != "" {
		filter = fmt.Sprintf("%s:%s", filterKey, filterValue)
	}

	key := &CacheKey{
		Type:   "docs:list",
		UserID: userID,
		Filter: filter,
		Limit:  limit,
	}
	return key.GenerateKey()
}

// DocumentKey создает ключ для отдельного документа
func DocumentKey(documentID string) string {
	key := &CacheKey{
		Type:       "docs:item",
		DocumentID: documentID,
	}
	return key.GenerateKey()
}

// UserKey создает ключ для пользователя
func UserKey(userID string) string {
	key := &CacheKey{
		Type:   "user",
		UserID: userID,
	}
	return key.GenerateKey()
}

// AccessKey создает ключ для прав доступа
func AccessKey(documentID, userID string) string {
	key := &CacheKey{
		Type:       "docs:access",
		DocumentID: documentID,
		UserID:     userID,
	}
	return key.GenerateKey()
}

// GetDocumentList получает список документов из кэша
func (cm *CacheManager) GetDocumentList(ctx context.Context, userID, filterKey, filterValue string, limit int) ([]interface{}, bool) {
	key := DocumentListKey(userID, filterKey, filterValue, limit)

	value, found := cm.cache.Get(ctx, key)
	if !found {
		return nil, false
	}

	if docs, ok := value.([]interface{}); ok {
		return docs, true
	}

	return nil, false
}

// SetDocumentList сохраняет список документов в кэш
func (cm *CacheManager) SetDocumentList(ctx context.Context, userID, filterKey, filterValue string, limit int, docs []interface{}) error {
	key := DocumentListKey(userID, filterKey, filterValue, limit)

	err := cm.cache.Set(ctx, key, docs, 5*time.Minute) // TTL 5 минут для списков
	if err != nil {
		return fmt.Errorf("failed to cache document list: %w", err)
	}

	return nil
}

// GetDocument получает документ из кэша
func (cm *CacheManager) GetDocument(ctx context.Context, documentID string) (interface{}, bool) {
	key := DocumentKey(documentID)

	value, found := cm.cache.Get(ctx, key)
	if !found {
		return nil, false
	}

	return value, true
}

// SetDocument сохраняет документ в кэш
func (cm *CacheManager) SetDocument(ctx context.Context, documentID string, doc interface{}) error {
	key := DocumentKey(documentID)

	err := cm.cache.Set(ctx, key, doc, 10*time.Minute) // TTL 10 минут для документов
	if err != nil {
		return fmt.Errorf("failed to cache document: %w", err)
	}

	return nil
}

// GetUser получает пользователя из кэша
func (cm *CacheManager) GetUser(ctx context.Context, userID string) (interface{}, bool) {
	key := UserKey(userID)

	value, found := cm.cache.Get(ctx, key)
	if !found {
		return nil, false
	}

	return value, true
}

// SetUser сохраняет пользователя в кэш
func (cm *CacheManager) SetUser(ctx context.Context, userID string, user interface{}) error {
	key := UserKey(userID)

	err := cm.cache.Set(ctx, key, user, 30*time.Minute) // TTL 30 минут для пользователей
	if err != nil {
		return fmt.Errorf("failed to cache user: %w", err)
	}

	return nil
}

// GetAccess получает права доступа из кэша
func (cm *CacheManager) GetAccess(ctx context.Context, documentID, userID string) (bool, bool) {
	key := AccessKey(documentID, userID)

	value, found := cm.cache.Get(ctx, key)
	if !found {
		return false, false
	}

	if hasAccess, ok := value.(bool); ok {
		return hasAccess, true
	}

	return false, false
}

// SetAccess сохраняет права доступа в кэш
func (cm *CacheManager) SetAccess(ctx context.Context, documentID, userID string, hasAccess bool) error {
	key := AccessKey(documentID, userID)

	err := cm.cache.Set(ctx, key, hasAccess, 15*time.Minute) // TTL 15 минут для прав доступа
	if err != nil {
		return fmt.Errorf("failed to cache access: %w", err)
	}

	return nil
}

// InvalidateDocument инвалидирует кэш для конкретного документа
func (cm *CacheManager) InvalidateDocument(ctx context.Context, documentID string) error {
	log.Printf("🗑️ Cache: Инвалидация кэша для документа %s", documentID)

	// Удаляем сам документ
	docKey := DocumentKey(documentID)
	if err := cm.cache.Delete(ctx, docKey); err != nil {
		return fmt.Errorf("failed to delete document from cache: %w", err)
	}

	// Инвалидируем все права доступа к документу
	accessPattern := fmt.Sprintf("docs:access:%s", documentID)
	if err := cm.cache.InvalidateByPattern(ctx, accessPattern); err != nil {
		return fmt.Errorf("failed to invalidate access cache: %w", err)
	}

	// Инвалидируем все списки документов (они могут содержать измененный документ)
	listPattern := "docs:list"
	if err := cm.cache.InvalidateByPattern(ctx, listPattern); err != nil {
		return fmt.Errorf("failed to invalidate document lists: %w", err)
	}

	return nil
}

// InvalidateUserDocuments инвалидирует кэш для всех документов пользователя
func (cm *CacheManager) InvalidateUserDocuments(ctx context.Context, userID string) error {
	log.Printf("🗑️ Cache: Инвалидация кэша для документов пользователя %s", userID)

	// Инвалидируем все списки документов пользователя
	listPattern := fmt.Sprintf("docs:list:%s", userID)
	if err := cm.cache.InvalidateByPattern(ctx, listPattern); err != nil {
		return fmt.Errorf("failed to invalidate user document lists: %w", err)
	}

	return nil
}

// InvalidateUser инвалидирует кэш для пользователя
func (cm *CacheManager) InvalidateUser(ctx context.Context, userID string) error {
	log.Printf("🗑️ Cache: Инвалидация кэша для пользователя %s", userID)

	// Удаляем пользователя
	userKey := UserKey(userID)
	if err := cm.cache.Delete(ctx, userKey); err != nil {
		return fmt.Errorf("failed to delete user from cache: %w", err)
	}

	// Инвалидируем все права доступа пользователя
	accessPattern := fmt.Sprintf("docs:access:%s", userID)
	if err := cm.cache.InvalidateByPattern(ctx, accessPattern); err != nil {
		return fmt.Errorf("failed to invalidate user access cache: %w", err)
	}

	// Инвалидируем списки документов пользователя
	return cm.InvalidateUserDocuments(ctx, userID)
}

// Clear очищает весь кэш
func (cm *CacheManager) Clear(ctx context.Context) error {
	log.Printf("🗑️ Cache: Очистка всего кэша")
	return cm.cache.Clear(ctx)
}
