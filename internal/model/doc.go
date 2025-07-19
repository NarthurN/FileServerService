package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
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
	log.Printf("JSONData.Scan: получено значение типа %T: %v", value, value)

	if value == nil {
		*j = nil
		log.Printf("JSONData.Scan: значение nil, устанавливаем пустое")
		return nil
	}

	// Попробуем разные типы
	switch v := value.(type) {
	case []byte:
		log.Printf("JSONData.Scan: обрабатываем []byte длиной %d", len(v))
		if len(v) == 0 {
			*j = make(JSONData)
			return nil
		}
		return json.Unmarshal(v, j)

	case string:
		log.Printf("JSONData.Scan: обрабатываем string: %s", v)
		if v == "" {
			*j = make(JSONData)
			return nil
		}
		return json.Unmarshal([]byte(v), j)

	case nil:
		*j = nil
		return nil

	default:
		log.Printf("JSONData.Scan: неизвестный тип %T", value)
		return fmt.Errorf("cannot scan %T into JSONData", value)
	}
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
	log.Printf("StringArray.Scan: получено значение типа %T: %v", value, value)

	if value == nil {
		*s = nil
		log.Printf("StringArray.Scan: значение nil, устанавливаем пустое")
		return nil
	}

	// Попробуем разные типы
	switch v := value.(type) {
	case []byte:
		log.Printf("StringArray.Scan: обрабатываем []byte длиной %d", len(v))
		if len(v) == 0 {
			*s = make(StringArray, 0)
			return nil
		}
		return json.Unmarshal(v, s)

	case string:
		log.Printf("StringArray.Scan: обрабатываем string: %s", v)
		if v == "" || v == "[]" {
			*s = make(StringArray, 0)
			return nil
		}
		return json.Unmarshal([]byte(v), s)

	case nil:
		*s = nil
		return nil

	default:
		log.Printf("StringArray.Scan: неизвестный тип %T", value)
		return fmt.Errorf("cannot scan %T into StringArray", value)
	}
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
