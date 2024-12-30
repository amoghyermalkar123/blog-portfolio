// internal/router/router.go
package router

import (
	"blog-portfolio/internal/config"
	"blog-portfolio/internal/handlers"
	"blog-portfolio/internal/logger"
	custommw "blog-portfolio/internal/middleware"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Router struct {
	chi.Router
	logger   *logger.Logger
	config   *config.Config
	handlers *handlers.Handlers
}

func New(logger *logger.Logger, config *config.Config, handlers *handlers.Handlers) *Router {
	r := chi.NewRouter()

	// Create custom middleware
	m := custommw.New(logger)

	// Log middleware setup
	logger.Info("Setting up middleware...")

	// Basic middleware
	r.Use(m.RequestLogger)
	// r.Use(m.Recover)
	r.Use(chimw.RealIP)
	r.Use(chimw.Timeout(60 * time.Second))

	// CORS middleware
	logger.Info("Configuring CORS with allowed origins:", config.Server.AllowOrigins)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{config.Server.AllowOrigins},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router := &Router{
		Router:   r,
		logger:   logger,
		config:   config,
		handlers: handlers,
	}

	logger.Info("Setting up routes...")
	router.setupRoutes()
	logger.Info("Router initialization complete")

	return router
}

// setupRoutes configures all application routes
func (router *Router) setupRoutes() {
	r := router.Router

	// Serve static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Authentication routes - these must come BEFORE the admin routes
	r.Group(func(r chi.Router) {
		r.Get("/login", router.handlers.Auth().ShowLogin())
		r.Post("/login", router.handlers.Auth().HandleLogin())
		r.Get("/logout", router.handlers.Auth().HandleLogout())
	})

	// Public routes
	r.Get("/", router.handlers.Home())
	r.Get("/blog", router.handlers.Posts().ListPosts())
	r.Get("/blog/{slug}", router.handlers.Posts().GetPost())

	// Admin routes - protected by RequireAuth middleware
	r.Route("/admin", func(r chi.Router) {
		r.Use(custommw.RequireAuth)

		// Dashboard
		r.Get("/dashboard", router.handlers.Admin().ShowDashboard())

		// Posts management
		r.Route("/posts", func(r chi.Router) {
			r.Get("/", router.handlers.Admin().ShowPosts())
			r.Get("/new/", router.handlers.Admin().ShowCreatePost())
			r.Get("/{id}", router.handlers.Admin().ShowEditPost())
			r.Post("/", router.handlers.Admin().HandleCreatePost())
			r.Put("/{id}", router.handlers.Admin().HandleUpdatePost())
			r.Delete("/{id}", router.handlers.Admin().HandleDeletePost())
		})
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/tags", func(r chi.Router) {
			r.Use(custommw.RequireAuth) // Protect these endpoints
			r.Post("/", router.handlers.Tags().CreateTag())
			r.Put("/{id}", router.handlers.Tags().UpdateTag())
			r.Delete("/{id}", router.handlers.Tags().DeleteTag())
		})
	})
}
