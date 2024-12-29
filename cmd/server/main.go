// cmd/server/main.go
package main

import (
	"blog-portfolio/internal/config"
	"blog-portfolio/internal/database"
	"blog-portfolio/internal/handlers"
	"blog-portfolio/internal/logger"
	"blog-portfolio/internal/repository"
	"blog-portfolio/internal/router"
	"blog-portfolio/internal/service"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize logger
	log := logger.New()

	// Load configuration
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	cfg, err := config.LoadConfig(env)
	if err != nil {
		log.Error("Failed to load configuration:", err)
		os.Exit(1)
	}

	// Initialize database
	db, err := database.New(log)
	if err != nil {
		log.Error("Failed to initialize database:", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations
	if err := db.RunMigrations(); err != nil {
		log.Error("Failed to run migrations:", err)
		os.Exit(1)
	}

	// Initialize repositories
	postRepo := repository.NewPostRepository(db.DB)

	// Initialize services
	postService := service.NewPostService(postRepo)

	// Initialize handlers
	h := handlers.New(log, postService)

	// Initialize router
	r := router.New(log, cfg, h)

	// Setup HTTP server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process lifecycle management
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Error("Graceful shutdown timed out... forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Error(err)
		}
		serverStopCtx()
	}()

	// Start the server
	log.Info("Server is running on http://localhost" + addr)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Error("Server failed to start:", err)
		os.Exit(1)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
