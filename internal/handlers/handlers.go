// internal/handlers/handlers.go
package handlers

import (
	"blog-portfolio/internal/logger"
	"blog-portfolio/internal/middleware"
	"blog-portfolio/internal/models"
	"blog-portfolio/internal/service"
	"blog-portfolio/web/layouts"
	"blog-portfolio/web/pages"
	"net/http"
)

type Handlers struct {
	logger      *logger.Logger
	posts       *PostHandlers
	auth        *AuthHandlers
	admin       *AdminHandlers
	postService *service.PostService
}

// New creates a new instance of Handlers
func New(logger *logger.Logger, postService *service.PostService, tagService *service.TagService) *Handlers {
	return &Handlers{
		logger:      logger,
		posts:       NewPostHandlers(postService, logger),
		auth:        NewAuthHandlers(logger),
		admin:       NewAdminHandlers(logger, postService, tagService), // Pass tagService here
		postService: postService,
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

// Home handles the home page
func (h *Handlers) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get latest posts - fixing the service call
		published := true
		latestPosts, err := h.postService.ListPosts(ctx, models.PostFilter{
			Published: &published,
			Limit:     3,
		})
		if err != nil {
			h.logger.Error("Error fetching latest posts:", err)
			latestPosts = []*models.Post{} // Empty slice if error
		}

		// Pass data to template
		if err := pages.Home(layouts.PageData{
			Title:       "Amogh's Eden",
			Description: "Welcome to my personal blog and portfolio",
			IsAdmin:     middleware.IsAdmin(r),
		}, latestPosts).Render(ctx, w); err != nil {
			h.logger.Error("Error rendering home page:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
