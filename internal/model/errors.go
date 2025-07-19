package model

import "errors"

// Стандартные ошибки репозитория
var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

// Бизнес-ошибки сервисного слоя
var (
	// Ошибки аутентификации
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidAdminToken  = errors.New("invalid admin token")

	// Ошибки валидации пользователя
	ErrLoginTooShort      = errors.New("login too short")
	ErrLoginTooLong       = errors.New("login too long")
	ErrLoginInvalidChars  = errors.New("login contains invalid characters")
	ErrLoginAlreadyExists = errors.New("login already exists")

	ErrPasswordTooShort    = errors.New("password too short")
	ErrPasswordTooLong     = errors.New("password too long")
	ErrPasswordNoUppercase = errors.New("password must contain uppercase letter")
	ErrPasswordNoLowercase = errors.New("password must contain lowercase letter")
	ErrPasswordNoDigit     = errors.New("password must contain digit")
	ErrPasswordNoSpecial   = errors.New("password must contain special character")

	// Ошибки валидации документов
	ErrDocumentNameEmpty    = errors.New("document name is empty")
	ErrDocumentNameTooLong  = errors.New("document name too long")
	ErrDocumentNameExists   = errors.New("document with this name already exists")
	ErrDocumentNoContent    = errors.New("document must have either file or JSON data")
	ErrDocumentInvalidGrant = errors.New("invalid grant user")

	// Ошибки прав доступа
	ErrAccessDenied      = errors.New("access denied")
	ErrOwnershipRequired = errors.New("only document owner can perform this action")

	// Общие ошибки валидации
	ErrRequired     = errors.New("required field is missing")
	ErrInvalidInput = errors.New("invalid input")
)

// BusinessError - структура для бизнес-ошибок с дополнительным контекстом
type BusinessError struct {
	Code    string
	Message string
	Cause   error
}

func (e BusinessError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e BusinessError) Unwrap() error {
	return e.Cause
}

// Конструкторы для бизнес-ошибок
func NewValidationError(message string, cause error) BusinessError {
	return BusinessError{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Cause:   cause,
	}
}

func NewAuthError(message string, cause error) BusinessError {
	return BusinessError{
		Code:    "AUTH_ERROR",
		Message: message,
		Cause:   cause,
	}
}

func NewAccessError(message string, cause error) BusinessError {
	return BusinessError{
		Code:    "ACCESS_ERROR",
		Message: message,
		Cause:   cause,
	}
}

func NewBusinessError(message string, cause error) BusinessError {
	return BusinessError{
		Code:    "BUSINESS_ERROR",
		Message: message,
		Cause:   cause,
	}
}
