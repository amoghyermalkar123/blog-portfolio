// internal/handlers/handlers.go
package handlers

import (
	"blog-portfolio/internal/logger"
	"blog-portfolio/internal/middleware"
	"blog-portfolio/internal/service"
	"blog-portfolio/web/layouts"
	"blog-portfolio/web/pages"
	"net/http"
)

type Handlers struct {
	logger *logger.Logger
	posts  *PostHandlers
	auth   *AuthHandlers
	admin  *AdminHandlers
	tags   *TagHandlers // New field for tag handlers
}

// New creates a new instance of Handlers
func New(logger *logger.Logger, postService *service.PostService, tagService *service.TagService) *Handlers {
	return &Handlers{
		logger: logger,
		posts:  NewPostHandlers(postService, logger),
		auth:   NewAuthHandlers(logger),
		admin:  NewAdminHandlers(logger, postService),
		tags:   NewTagHandlers(logger, tagService), // Initialize tag handlers
	}
}

// Posts returns the post handlers
func (h *Handlers) Posts() *PostHandlers {
	return h.posts
}

// Auth returns the auth handlers
func (h *Handlers) Auth() *AuthHandlers {
	return h.auth
}

// Admin returns the admin handlers
func (h *Handlers) Admin() *AdminHandlers {
	return h.admin
}

// Add a new method to access tag handlers
func (h *Handlers) Tags() *TagHandlers {
	return h.tags
}

// Home handles the home page
func (h *Handlers) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAdmin := middleware.IsAdmin(r)

		if err := pages.Home(layouts.PageData{
			Title:       "Home | Blog & Portfolio",
			Description: "Welcome to my personal blog and portfolio",
			IsAdmin:     isAdmin,
		}).Render(r.Context(), w); err != nil {
			h.logger.Error("Error rendering home page:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
