package cache

import (
	"context"
	"sync"
	"time"
)

// CacheItem представляет элемент кэша
type CacheItem struct {
	Key       string
	Value     interface{}
	ExpiresAt time.Time
	CreatedAt time.Time
	Next      *CacheItem
	Prev      *CacheItem
}

// IsExpired проверяет, истек ли срок действия элемента
func (item *CacheItem) IsExpired() bool {
	return !item.ExpiresAt.IsZero() && time.Now().After(item.ExpiresAt)
}

// Cache интерфейс для кэша
type Cache interface {
	// Get получает значение по ключу
	Get(ctx context.Context, key string) (interface{}, bool)

	// Set устанавливает значение с TTL
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Delete удаляет элемент по ключу
	Delete(ctx context.Context, key string) error

	// InvalidateByPattern удаляет все элементы, соответствующие паттерну
	InvalidateByPattern(ctx context.Context, pattern string) error

	// Clear очищает весь кэш
	Clear(ctx context.Context) error
}

// LRUCache реализация LRU кэша
type LRUCache struct {
	mu       sync.RWMutex
	capacity int64
	items    map[string]*CacheItem
	head     *CacheItem
	tail     *CacheItem
}

// NewLRUCache создает новый LRU кэш
func NewLRUCache(capacity int64) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*CacheItem),
	}
}

// Get получает значение по ключу
func (c *LRUCache) Get(ctx context.Context, key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Проверяем истечение срока действия
	if item.IsExpired() {
		c.removeItem(item)
		return nil, false
	}

	// Перемещаем элемент в начало списка (LRU)
	c.moveToFront(item)

	return item.Value, true
}

// Set устанавливает значение с TTL
func (c *LRUCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Если элемент уже существует, обновляем его
	if existingItem, exists := c.items[key]; exists {
		existingItem.Value = value
		existingItem.ExpiresAt = time.Now().Add(ttl)
		c.moveToFront(existingItem)
		return nil
	}

	// Создаем новый элемент
	item := &CacheItem{
		Key:       key,
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
	}

	// Если кэш полон, удаляем самый старый элемент
	if int64(len(c.items)) >= c.capacity {
		c.evictOldest()
	}

	// Добавляем новый элемент
	c.items[key] = item
	c.addToFront(item)

	return nil
}

// Delete удаляет элемент по ключу
func (c *LRUCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, exists := c.items[key]; exists {
		c.removeItem(item)
	}
	return nil
}

// InvalidateByPattern удаляет все элементы, соответствующие паттерну
func (c *LRUCache) InvalidateByPattern(ctx context.Context, pattern string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Удаляем все элементы, содержащие паттерн
	for key, item := range c.items {
		if containsPattern(key, pattern) {
			c.removeItem(item)
		}
	}
	return nil
}

// Clear очищает весь кэш
func (c *LRUCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*CacheItem)
	c.head = nil
	c.tail = nil
	return nil
}

// moveToFront перемещает элемент в начало списка
func (c *LRUCache) moveToFront(item *CacheItem) {
	if c.head == item {
		return
	}

	// Удаляем из текущей позиции
	c.removeFromList(item)
	// Добавляем в начало
	c.addToFront(item)
}

// addToFront добавляет элемент в начало списка
func (c *LRUCache) addToFront(item *CacheItem) {
	if c.head == nil {
		c.head = item
		c.tail = item
		return
	}

	item.Next = c.head
	item.Prev = nil
	c.head.Prev = item
	c.head = item
}

// removeFromList удаляет элемент из списка
func (c *LRUCache) removeFromList(item *CacheItem) {
	if item.Prev != nil {
		item.Prev.Next = item.Next
	} else {
		c.head = item.Next
	}

	if item.Next != nil {
		item.Next.Prev = item.Prev
	} else {
		c.tail = item.Prev
	}

	item.Next = nil
	item.Prev = nil
}

// removeItem удаляет элемент из кэша
func (c *LRUCache) removeItem(item *CacheItem) {
	c.removeFromList(item)
	delete(c.items, item.Key)
}

// evictOldest удаляет самый старый элемент
func (c *LRUCache) evictOldest() {
	if c.tail != nil {
		c.removeItem(c.tail)
	}
}

// containsPattern проверяет, содержит ли строка паттерн
func containsPattern(str, pattern string) bool {
	return len(pattern) > 0 && len(str) >= len(pattern) &&
		(str == pattern ||
			(len(str) > len(pattern) &&
				(str[:len(pattern)] == pattern ||
					str[len(str)-len(pattern):] == pattern ||
					containsSubstring(str, pattern))))
}

// containsSubstring проверяет, содержит ли строка подстроку
func containsSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
