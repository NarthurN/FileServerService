# FileServerService

## 1. Что за проект и зачем нужен

**FileServerService** — это микросервис для сохранения и раздачи электронных документов с REST API. Сервис предоставляет безопасную систему управления документами с аутентификацией пользователей, поддержкой различных типов файлов и метаданных.

**Основные возможности:**
- Регистрация и аутентификация пользователей
- Загрузка документов (файлы и JSON данные)
- Получение списка документов с фильтрацией
- Скачивание документов по ID
- Удаление документов
- Управление правами доступа к документам

## 2. Особенности, технологии и библиотеки

### Основные технологии:
- **Go 1.23.4** — основной язык разработки
- **PostgreSQL 15** — реляционная база данных
- **Docker & Docker Compose** — контейнеризация и оркестрация

### Ключевые библиотеки и их назначение:

| Библиотека | Назначение |
|------------|------------|
| `github.com/ogen-go/ogen` | Генерация кода из OpenAPI спецификации |
| `github.com/go-chi/chi/v5` | HTTP роутер с middleware |
| `github.com/jackc/pgx/v5` | Драйвер PostgreSQL с пулом соединений |
| `github.com/Masterminds/squirrel` | SQL query builder |
| `github.com/google/uuid` | Генерация UUID |
| `github.com/joho/godotenv` | Загрузка переменных окружения |

### Особенности архитектуры:
- **Clean Architecture** — разделение на слои (API, Service, Repository)
- **OpenAPI/ogen** — автоматическая генерация кода из спецификации
- **In-memory кэш** — для повышения производительности
- **JWT токены** — для аутентификации пользователей
- **Миграции БД** — автоматическое управление схемой

## 3. Как запустить приложение

### Быстрый запуск с Docker

1. **Клонируйте репозиторий:**
   ```bash
   git clone <repository-url>
   cd FileServerService
   ```

2. **Запустите приложение:**
   ```bash
   docker compose up -d
   ```

3. **Проверьте статус:**
   ```bash
   docker compose ps
   ```

### Доступ к приложению

После запуска приложение доступно по следующим адресам:

- **API**: http://localhost:8080


### Локальная разработка

1. **Установите Go 1.23.4+**

2. **Запустите PostgreSQL:**
   ```bash
   docker compose up -d postgres
   ```

3. **Запустите приложение:**
   ```bash
   go run ./cmd/server
   ```

### Переменные окружения

```env
# База данных
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=docs_server
DB_SSL_MODE=disable

# Сервер
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Аутентификация
ADMIN_TOKEN=super-secret-admin-token-for-user-registration-2024
TOKEN_LIFETIME_HOURS=24
JWT_SECRET=your-very-long-and-secure-jwt-secret-key-here
```

## 4. API Endpoints

| Метод | Endpoint | Описание | Аутентификация |
|-------|----------|----------|----------------|
| `POST` | `/api/register` | Регистрация пользователя | Admin Token |
| `POST` | `/api/auth` | Авторизация пользователя | - |
| `DELETE` | `/api/auth/{token}` | Выход из системы | Token |
| `GET` | `/api/docs` | Список документов | Token |
| `POST` | `/api/docs` | Создание документа | Token |
| `GET` | `/api/docs/{id}` | Получение документа | Token |
| `DELETE` | `/api/docs/{id}` | Удаление документа | Token |

### Примеры curl запросов

#### Регистрация пользователя
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer super-secret-admin-token-for-user-registration-2024" \
  -d '{
    "login": "user@example.com",
    "password": "password123"
  }'
```

#### Авторизация пользователя
```bash
curl -X POST http://localhost:8080/api/auth \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "login=user@example.com&password=password123"
```

#### Создание документа (файл)
```bash
curl -X POST http://localhost:8080/api/docs \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/document.pdf" \
  -F "name=My Document" \
  -F "description=Document description"
```

#### Создание документа (JSON)
```bash
curl -X POST http://localhost:8080/api/docs \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "JSON Document",
    "description": "Document with JSON data",
    "data": {"key": "value", "number": 42}
  }'
```

#### Получение списка документов
```bash
curl -X GET "http://localhost:8080/api/docs?token=YOUR_TOKEN&limit=10"
```

#### Получение документа по ID
```bash
curl -X GET http://localhost:8080/api/docs/DOCUMENT_ID \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### Удаление документа
```bash
curl -X DELETE http://localhost:8080/api/docs/DOCUMENT_ID \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### Выход из системы
```bash
curl -X DELETE http://localhost:8080/api/auth/YOUR_TOKEN
```

## 5. Архитектура проекта

Проект построен по принципам **Clean Architecture** с четким разделением ответственности:

### Структура папок:

```
FileServerService/
├── cmd/server/           # Точка входа в приложение
├── internal/             # Внутренняя логика (не экспортируется)
│   ├── api/v1/          # HTTP handlers и валидация
│   ├── cache/           # In-memory кэш для производительности
│   ├── config/          # Конфигурация приложения
│   ├── database/        # Слой работы с БД и миграции
│   ├── model/           # Доменные модели и ошибки
│   ├── repository/      # Слой доступа к данным
│   └── service/         # Бизнес-логика и use cases
├── pkg/                 # Публичные пакеты
│   ├── generated/       # Автогенерированный код из OpenAPI
│   └── openapi/         # OpenAPI спецификации
└── db/                  # Инициализация БД
```

### Назначение каждой папки:

| Папка | Назначение |
|-------|------------|
| `cmd/server/` | Точка входа в приложение, инициализация зависимостей |
| `internal/api/v1/` | HTTP handlers, валидация запросов, обработка ошибок |
| `internal/cache/` | In-memory кэш для кэширования часто запрашиваемых данных |
| `internal/config/` | Загрузка и валидация конфигурации из переменных окружения |
| `internal/database/` | Подключение к БД, пул соединений, миграции |
| `internal/model/` | Доменные модели (User, Document, Token), кастомные ошибки |
| `internal/repository/` | Слой доступа к данным, SQL запросы, CRUD операции |
| `internal/service/` | Бизнес-логика, аутентификация, валидация прав доступа |
| `pkg/generated/` | Автогенерированный код из OpenAPI спецификации |
| `pkg/openapi/` | OpenAPI спецификации для генерации кода и документации |

### Слои архитектуры:

1. **API Layer** (`internal/api/`) — обработка HTTP запросов
2. **Service Layer** (`internal/service/`) — бизнес-логика
3. **Repository Layer** (`internal/repository/`) — доступ к данным
4. **Database Layer** (`internal/database/`) — работа с БД

### Поток данных:
```
HTTP Request → API Handler → Service → Repository → Database
                ↓
HTTP Response ← API Handler ← Service ← Repository ← Database
```

### Управление контейнерами

```bash
# Остановить приложение
docker compose down

# Просмотр логов
docker compose logs -f fileserver

# Пересборка образа
docker compose build

# Перезапуск
docker compose restart
```
