package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rentalflow/inventory-service/internal/config"
	"github.com/rentalflow/inventory-service/internal/handler"
	"github.com/rentalflow/inventory-service/internal/repository"
	"github.com/rentalflow/inventory-service/internal/service"
	"github.com/rentalflow/rentalflow/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Init(cfg.ServiceName, cfg.LogLevel)
	log := logger.NewLogger("main")

	log.Info().Msg("Starting Inventory Service...")

	// Connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.Database.DSN())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer pool.Close()

	log.Info().Str("host", cfg.Database.Host).Msg("Connected to database")

	// Initialize repositories
	itemRepo := repository.NewPostgresItemRepository(pool)
	availabilityRepo := repository.NewPostgresAvailabilityRepository(pool)
	maintenanceRepo := repository.NewPostgresMaintenanceRepository(pool)

	// Initialize service
	inventoryService := service.NewInventoryService(itemRepo, availabilityRepo, maintenanceRepo)

	// Initialize HTTP handler
	httpHandler := handler.NewHTTPHandler(inventoryService)

	// Start HTTP server
	httpAddr := fmt.Sprintf(":%d", cfg.HTTPPort)
	mux := http.NewServeMux()
	httpHandler.RegisterRoutes(mux)

	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	go func() {
		log.Info().Str("addr", httpAddr).Msg("HTTP API server listening")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("HTTP server failed")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Graceful shutdown
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown failed")
	}

	log.Info().Msg("Server stopped")
}
