package validate

import (
	"fmt"
	"mime"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/NarthurN/FileServerService/internal/model"
)

// FileValidator - валидатор файлов
type FileValidator struct {
	maxFileSize  int64
	allowedTypes map[string]bool
	blockedTypes map[string]bool
}

func NewFileValidator() *FileValidator {
	return &FileValidator{
		maxFileSize: 100 * 1024 * 1024, // 100MB по умолчанию
		allowedTypes: map[string]bool{
			"image/jpeg":               true,
			"image/png":                true,
			"image/gif":                true,
			"application/pdf":          true,
			"text/plain":               true,
			"application/json":         true,
			"application/octet-stream": true,
		},
		blockedTypes: map[string]bool{
			"application/x-executable": true,
			"application/x-msdownload": true,
		},
	}
}

func (fv *FileValidator) ValidateFile(filename, mimeType string, size int64) error {
	// Проверка размера
	if size > fv.maxFileSize {
		return fmt.Errorf("file too large: %d bytes (max %d)", size, fv.maxFileSize)
	}

	// Проверка MIME типа
	if fv.blockedTypes[mimeType] {
		return fmt.Errorf("file type not allowed: %s", mimeType)
	}

	// Если есть whitelist, проверяем его
	if len(fv.allowedTypes) > 0 && !fv.allowedTypes[mimeType] {
		return fmt.Errorf("file type not supported: %s", mimeType)
	}

	// Проверка расширения файла
	ext := strings.ToLower(filepath.Ext(filename))
	expectedMime := mime.TypeByExtension(ext)
	if expectedMime != "" && expectedMime != mimeType {
		return fmt.Errorf("MIME type mismatch: expected %s for %s, got %s", expectedMime, ext, mimeType)
	}

	return nil
}

// DocumentNameValidator - валидатор имен документов
type DocumentNameValidator struct {
	maxLength      int
	forbiddenChars *regexp.Regexp
	reservedNames  map[string]bool
}

func NewDocumentNameValidator() *DocumentNameValidator {
	return &DocumentNameValidator{
		maxLength:      255,
		forbiddenChars: regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`),
		reservedNames: map[string]bool{
			"con": true, "prn": true, "aux": true, "nul": true,
			"com1": true, "com2": true, "com3": true, "com4": true,
			"lpt1": true, "lpt2": true, "lpt3": true, "lpt4": true,
		},
	}
}

func (dnv *DocumentNameValidator) ValidateName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("document name cannot be empty")
	}

	if len(name) > dnv.maxLength {
		return fmt.Errorf("document name too long: %d characters (max %d)", len(name), dnv.maxLength)
	}

	if dnv.forbiddenChars.MatchString(name) {
		return fmt.Errorf("document name contains forbidden characters")
	}

	// Проверка на зарезервированные имена (Windows)
	baseName := strings.ToLower(strings.TrimSuffix(name, filepath.Ext(name)))
	if dnv.reservedNames[baseName] {
		return fmt.Errorf("document name is reserved: %s", name)
	}

	return nil
}

// AccessManager - управление правами доступа
type AccessManager struct{}

func NewAccessManager() *AccessManager {
	return &AccessManager{}
}

func (am *AccessManager) CanAccessDocument(doc model.Document, userID, userLogin string) bool {
	// Владелец всегда имеет доступ
	if doc.UserID == userID {
		return true
	}

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

func (am *AccessManager) CanModifyDocument(doc model.Document, userID string) bool {
	// Только владелец может изменять документ
	return doc.UserID == userID
}

func (am *AccessManager) CanDeleteDocument(doc model.Document, userID string) bool {
	// Только владелец может удалять документ
	return doc.UserID == userID
}

// RateLimiter - простой rate limiter для операций
type RateLimiter struct {
	attempts    map[string][]time.Time
	maxAttempts int
	window      time.Duration
}

func NewRateLimiter(maxAttempts int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		attempts:    make(map[string][]time.Time),
		maxAttempts: maxAttempts,
		window:      window,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	now := time.Now()

	// Очищаем старые попытки
	attempts := rl.attempts[key]
	validAttempts := make([]time.Time, 0, len(attempts))

	for _, attempt := range attempts {
		if now.Sub(attempt) < rl.window {
			validAttempts = append(validAttempts, attempt)
		}
	}

	// Проверяем лимит
	if len(validAttempts) >= rl.maxAttempts {
		rl.attempts[key] = validAttempts
		return false
	}

	// Добавляем новую попытку
	validAttempts = append(validAttempts, now)
	rl.attempts[key] = validAttempts

	return true
}

// SecurityUtils - утилиты безопасности
type SecurityUtils struct{}

func NewSecurityUtils() *SecurityUtils {
	return &SecurityUtils{}
}

func (su *SecurityUtils) SanitizeInput(input string) string {
	// Удаляем потенциально опасные символы
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "\x00", "") // null bytes
	input = strings.ReplaceAll(input, "\r", "")   // carriage returns

	return input
}

func (su *SecurityUtils) IsValidUUID(uuid string) bool {
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidPattern.MatchString(strings.ToLower(uuid))
}

func (su *SecurityUtils) MaskSensitiveData(data string) string {
	if len(data) <= 4 {
		return strings.Repeat("*", len(data))
	}

	return data[:2] + strings.Repeat("*", len(data)-4) + data[len(data)-2:]
}
