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
	"github.com/rentalflow/rentalflow/pkg/logger"
	"github.com/rentalflow/review-service/internal/config"
	"github.com/rentalflow/review-service/internal/handler"
	"github.com/rentalflow/review-service/internal/repository"
	"github.com/rentalflow/review-service/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger.Init(cfg.ServiceName, cfg.LogLevel)
	log := logger.NewLogger("main")
	log.Info().Msg("Starting Review Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.Database.DSN())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer pool.Close()

	log.Info().Str("host", cfg.Database.Host).Msg("Connected to database")

	// Run migration
	migrationSQL := `
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    DO $$ BEGIN
        CREATE TYPE review_type AS ENUM ('renter_to_owner', 'owner_to_renter', 'renter_to_item');
    EXCEPTION
        WHEN duplicate_object THEN null;
    END $$;

    CREATE TABLE IF NOT EXISTS reviews (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        booking_id UUID NOT NULL,
        reviewer_id UUID NOT NULL,
        target_user_id UUID,
        target_item_id UUID,
        review_type review_type NOT NULL,
        rating DECIMAL(2, 1) CHECK (rating >= 1.0 AND rating <= 5.0),
        comment TEXT NOT NULL,
        is_verified BOOLEAN DEFAULT FALSE,
        is_visible BOOLEAN DEFAULT TRUE,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );

    CREATE INDEX IF NOT EXISTS idx_reviews_booking ON reviews(booking_id);
    CREATE INDEX IF NOT EXISTS idx_reviews_reviewer ON reviews(reviewer_id);
    CREATE INDEX IF NOT EXISTS idx_reviews_target_user ON reviews(target_user_id);
    CREATE INDEX IF NOT EXISTS idx_reviews_target_item ON reviews(target_item_id);
    CREATE INDEX IF NOT EXISTS idx_reviews_type ON reviews(review_type);
    CREATE INDEX IF NOT EXISTS idx_reviews_visible ON reviews(is_visible);
    `
	_, err = pool.Exec(ctx, migrationSQL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run migration")
	}
	log.Info().Msg("Migration applied successfully")

	reviewRepo := repository.NewPostgresReviewRepository(pool)
	reviewService := service.NewReviewService(reviewRepo)
	httpHandler := handler.NewHTTPHandler(reviewService)

	httpAddr := fmt.Sprintf(":%d", cfg.HTTPPort)
	mux := http.NewServeMux()
	httpHandler.RegisterRoutes(mux)

	httpServer := &http.Server{Addr: httpAddr, Handler: mux}

	go func() {
		log.Info().Str("addr", httpAddr).Msg("HTTP API server listening")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("HTTP server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown failed")
	}

	log.Info().Msg("Server stopped")
}
