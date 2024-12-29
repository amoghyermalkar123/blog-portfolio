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

	// Basic middleware
	r.Use(m.RequestLogger)
	r.Use(m.Recover)
	r.Use(chimw.RealIP)
	r.Use(chimw.Timeout(60 * time.Second))

	// CORS middleware
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

	// Setup routes
	router.setupRoutes()

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

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Add API routes here
	})

	// Web routes
	r.Group(func(r chi.Router) {
		r.Get("/", router.handlers.Home())

		// Blog routes
		r.Route("/blog", func(r chi.Router) {
			r.Get("/", router.handlers.Posts().ListPosts())
			r.Get("/{slug}", router.handlers.Posts().GetPost())
		})

		// Portfolio routes
		r.Route("/portfolio", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Portfolio items will be listed here"))
			})
			r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
				id := chi.URLParam(r, "id")
				w.Write([]byte("Portfolio item: " + id))
			})
		})

		// Admin routes
		r.Route("/admin", func(r chi.Router) {
			// r.Use(custommw.RequireAuth)

			// Dashboard
			r.Get("/dashboard", router.handlers.Admin().ShowDashboard())

			// Post Management
			r.Route("/posts", func(r chi.Router) {
				r.Get("/", router.handlers.Admin().ShowPosts())
				r.Get("/new/", router.handlers.Admin().ShowCreatePost())
			})
		})
	})
}
