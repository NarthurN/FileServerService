package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Document - репозиторная модель документа
type Document struct {
	ID        string      `db:"id" json:"id"`                    // ID документа
	UserID    string      `db:"user_id" json:"-"`                // ID пользователя
	Name      string      `db:"name" json:"name"`                // Название документа
	MimeType  string      `db:"mime_type" json:"mime"`           // MIME-тип документа
	FilePath  string      `db:"file_path" json:"-"`              // Путь к файлу
	IsFile    bool        `db:"is_file" json:"file"`             // Флаг, является ли файл
	IsPublic  bool        `db:"is_public" json:"public"`         // Флаг, является ли документ публичным
	JSONData  JSONData    `db:"json_data" json:"json,omitempty"` // JSON данные документа
	Grants    StringArray `db:"grants" json:"grant"`             // Массив логинов с доступом
	CreatedAt time.Time   `db:"created_at" json:"created"`       // Дата создания документа
	UpdatedAt time.Time   `db:"updated_at" json:"-"`             // Дата обновления документа
}

// JSONData - тип для хранения JSON данных
type JSONData map[string]any

// Value реализует интерфейс driver.Valuer для записи в БД
func (j JSONData) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan реализует интерфейс sql.Scanner для чтения из БД
func (j *JSONData) Scan(value any) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("cannot scan non-[]byte value into JSONData")
	}

	return json.Unmarshal(bytes, j)
}

// StringArray - тип для хранения массива строк
type StringArray []string

// Value реализует интерфейс driver.Valuer для записи в БД
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Scan реализует интерфейс sql.Scanner для чтения из БД
func (s *StringArray) Scan(value any) error {
	if value == nil {
		*s = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("cannot scan non-[]byte value into StringArray")
	}

	return json.Unmarshal(bytes, s)
}

// MarshalJSON кастомная сериализация для времени в нужном формате
func (d Document) MarshalJSON() ([]byte, error) {
	type Alias Document

	return json.Marshal(&struct {
		*Alias
		Created string `json:"created"`
	}{
		Alias:   (*Alias)(&d),
		Created: d.CreatedAt.Format("2006-01-02 15:04:05"), // Формат из задания
	})
}
