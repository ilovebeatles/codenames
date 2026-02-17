package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"codenames/internal/config"
	"codenames/internal/game"
	"codenames/internal/handler"
	"codenames/internal/hub"
	"codenames/internal/storage"
)

func main() {
	cfg := config.Load()

	// Run migrations
	migrationsPath := "migrations"
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		migrationsPath = "backend/migrations"
	}
	if err := storage.RunMigrations(cfg.DatabaseURL, migrationsPath); err != nil {
		log.Printf("warning: migrations: %v", err)
	}

	// Connect to DB
	ctx := context.Background()
	pool, err := storage.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Init repos
	roomRepo := storage.NewRoomRepo(pool)
	playerRepo := storage.NewPlayerRepo(pool)
	gameRepo := storage.NewGameRepo(pool)

	// Init engine
	engine := game.NewEngine(gameRepo, playerRepo)

	// Init hub
	h := hub.NewHub(roomRepo, playerRepo, gameRepo, engine)
	go h.Run()

	// Init handlers
	roomHandler := handler.NewRoomHandler(roomRepo, playerRepo, gameRepo)
	playerHandler := handler.NewPlayerHandler(playerRepo)
	wsHandler := handler.NewWSHandler(h, playerRepo)

	// Init router
	r := handler.NewRouter(roomHandler, playerHandler, wsHandler)

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("server starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Println("server stopped")
}
