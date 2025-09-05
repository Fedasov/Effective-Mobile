package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fedasov/Effective-Mobile/internal/config"
	"github.com/Fedasov/Effective-Mobile/internal/handler"
	"github.com/Fedasov/Effective-Mobile/internal/middleware"
	"github.com/Fedasov/Effective-Mobile/internal/repository"
	"github.com/Fedasov/Effective-Mobile/internal/service"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	_ "github.com/Fedasov/Effective-Mobile/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Subscription Service API
// @version 1.0
// @description API для управления онлайн-подписками пользователей

// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.Load()

	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	router := setupRouter(subscriptionHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// initDB инициализирует подключение к PostgreSQL
func initDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return db, nil
}

func setupRouter(subscriptionHandler *handler.SubscriptionHandler) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)

	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/subscriptions", subscriptionHandler.Create).Methods("POST")
	api.HandleFunc("/subscriptions/{id}", subscriptionHandler.GetByID).Methods("GET")
	api.HandleFunc("/subscriptions/{id}", subscriptionHandler.Update).Methods("PUT")
	api.HandleFunc("/subscriptions/{id}", subscriptionHandler.Delete).Methods("DELETE")
	api.HandleFunc("/subscriptions", subscriptionHandler.List).Methods("GET")
	api.HandleFunc("/subscriptions/total-cost", subscriptionHandler.GetTotalCost).Methods("POST")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	}).Methods("GET")

	return router
}
