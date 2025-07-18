package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	fileserverAPI "github.com/NarthurN/FileServerService/internal/api/docs"
	"github.com/NarthurN/FileServerService/internal/config"
	"github.com/NarthurN/FileServerService/internal/database"
	"github.com/NarthurN/FileServerService/internal/database/migrator"
	fileserverCompositeRepo "github.com/NarthurN/FileServerService/internal/repository"
	fileserverService "github.com/NarthurN/FileServerService/internal/service/docs"
	fileserverV1 "github.com/NarthurN/FileServerService/pkg/generated/api/fileserver/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Контекст для миграций
	ctx := context.Background()

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Printf("🚨 ошибка загрузки конфигурации: %v", err)
		return
	}

	// Создание SQL соединения
	sqlDB, err := database.NewSQLDB(cfg.Database)
    if err != nil {
        log.Fatal("🚨 ошибка создания SQL соединения:", err)
    }
    defer sqlDB.Close()

	// Применение миграций
	migrator := migrator.NewMigrator(sqlDB)
    if err := migrator.Up(ctx); err != nil {
        log.Fatal("🚨 ошибка применения миграций:", err)
    }

	// Создание пула соединений
	pool, err := database.NewPool(cfg.Database)
	if err != nil {
		log.Printf("🚨 ошибка создания пула соединений: %v", err)
		return
	}

	// Создание репозитория
	repo := fileserverCompositeRepo.NewCompositeRepository(pool)

	// Создание сервиса
	service := fileserverService.NewService(repo)

	// Создание API
	api := fileserverAPI.NewAPI(service)

	// Создание сервера
	fileServer, err := fileserverV1.NewServer(api)
	if err != nil {
		log.Printf("🚨 ошибка создания сервера: %v", err)
		return
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Mount("/api", fileServer)

    // Статические файлы для загруженных документов
    r.Handle("/uploads/*", http.StripPrefix("/uploads", http.FileServer(http.Dir("./bin/storage"))))

    // Swagger UI
	swaggerFS := http.FileServer(http.Dir("./pkg/openapi/bundles"))
    r.Handle("/swagger-ui.html", swaggerFS)
    r.Handle("/docs.openapi.bundle.yaml", swaggerFS)
    r.Handle("/swagger/*", http.StripPrefix("/swagger", swaggerFS))

    // Главная страница
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
    })

    server := &http.Server{
        Addr:              net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.Port)),
        Handler:           r, // Используем chi router напрямую
        ReadHeaderTimeout: 10 * time.Second,
        WriteTimeout:      60 * time.Second,  // Для файлов
        IdleTimeout:       120 * time.Second,
        MaxHeaderBytes:    32 << 20, // 32MB для файлов
    }

	go func() {
		log.Printf("🚀 HTTP сервер запущен на порту %s", strconv.Itoa(cfg.Server.Port))
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("🚨 ошибка запуска HTTP сервера: %v", err)
		}
	}()

	// Gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("👋 HTTP сервер завершает работу...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("🚨 ошибка завершения работы HTTP сервера: %v", err)
	}

	log.Println("👋 HTTP сервер завершен")
}
