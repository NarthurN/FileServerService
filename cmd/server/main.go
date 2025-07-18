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
	fileserverRepo "github.com/NarthurN/FileServerService/internal/repository/doc"
	fileserverService "github.com/NarthurN/FileServerService/internal/service/docs"
	fileserverV1 "github.com/NarthurN/FileServerService/pkg/generated/api/fileserver/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("🚨 ошибка загрузки конфигурации: %v", err)
		return
	}

	db, err := database.NewPool(cfg.Database)
	if err != nil {
		log.Printf("🚨 ошибка создания пула соединений: %v", err)
		return
	}

	repo := fileserverRepo.NewRepository(db)
	service := fileserverService.NewService(repo)
	api := fileserverAPI.NewAPI(service)

	fileServer, err := fileserverV1.NewServer(api, nil)
	if err != nil {
		log.Printf("🚨 ошибка создания сервера: %v", err)
		return
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/api", fileServer)

	// Файл сервер для документации
	swaggerServer := http.FileServer(http.Dir("./pkg/openapi/bundles"))

	httpMux := http.NewServeMux()
	httpMux.Handle("/api/v1/", http.StripPrefix("/api/v1", r))

	httpMux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
			return
		}
		swaggerServer.ServeHTTP(w, r)
	}))

	server := &http.Server{
		Addr:              net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.Port)),
		Handler:           httpMux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       15 * time.Second,
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
