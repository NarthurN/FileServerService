# Многоэтапный Dockerfile для Go приложения
FROM golang:1.23.4-alpine AS builder

# Установка необходимых пакетов для сборки
RUN apk add --no-cache git ca-certificates tzdata

# Установка рабочей директории
WORKDIR /app

# Копирование файлов зависимостей
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Финальный этап - минимальный образ
FROM alpine:latest

# Установка необходимых пакетов для runtime
RUN apk --no-cache add ca-certificates tzdata

# Создание пользователя для безопасности
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Создание директорий для приложения
WORKDIR /app

# Создание директории для хранения файлов
RUN mkdir -p /app/bin/storage && \
    chown -R appuser:appgroup /app

# Копирование скомпилированного приложения
COPY --from=builder /app/main .

# Копирование статических файлов (если есть)
COPY --from=builder /app/pkg/openapi/bundles ./pkg/openapi/bundles

# Переключение на непривилегированного пользователя
USER appuser

# Открытие порта
EXPOSE 8080

# Проверка здоровья
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/docs || exit 1

# Запуск приложения
CMD ["./main"]
