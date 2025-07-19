# Настройки базы данных
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=docs_server
DB_SSL_MODE=disable

# Docker PostgreSQL настройки
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=docs_server

# Настройки сервера
SERVER_HOST=localhost
SERVER_PORT=8080

# Goose настройки
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://postgres:postgres@localhost:5432/docs_server?sslmode=disable
GOOSE_MIGRATION_DIR=db/migrations

# Настройки аутентификации
ADMIN_TOKEN=super-secret-admin-token-for-user-registration-2024
TOKEN_LIFETIME_HOURS=24
JWT_SECRET=your-very-long-and-secure-jwt-secret-key-here
