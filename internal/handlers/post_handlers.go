// internal/handlers/post_handlers.go
package handlers

import (
	"blog-portfolio/internal/logger"
	"blog-portfolio/internal/models"
	"blog-portfolio/internal/service"
	"blog-portfolio/web/pages"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PostHandlers struct {
	service *service.PostService
	logger  *logger.Logger
}

func NewPostHandlers(service *service.PostService, logger *logger.Logger) *PostHandlers {
	return &PostHandlers{
		service: service,
		logger:  logger,
	}
}

// ListPosts handles the blog listing page
// internal/handlers/post_handlers.go

func (h *PostHandlers) ListPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse query parameters
		tag := r.URL.Query().Get("tag")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}

		// Set up pagination
		limit := 10
		offset := (page - 1) * limit

		// Important: Set published filter to true for public blog page
		published := true

		// Create filter
		filter := models.PostFilter{
			Tag:       tag,
			Published: &published, // Make sure we're passing the address of published
			Limit:     limit,
			Offset:    offset,
		}

		// Get posts
		posts, err := h.service.ListPosts(ctx, filter)
		if err != nil {
			h.logger.Error("Error listing posts:", err)
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}

		// Handle different response types
		switch {
		case r.Header.Get("HX-Request") == "true":
			err = pages.BlogPostList(posts).Render(ctx, w)
		case r.Header.Get("Accept") == "application/json":
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(posts); err != nil {
				h.logger.Error("Error encoding posts:", err)
				http.Error(w, "Error encoding response", http.StatusInternalServerError)
			}
			return
		default:
			err = pages.Blog(posts, page, tag).Render(ctx, w)
		}

		if err != nil {
			h.logger.Error("Error rendering blog page:", err)
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
		}
	}
}

// GetPost handles individual blog post pages
func (h *PostHandlers) GetPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		slug := chi.URLParam(r, "slug")

		post, err := h.service.GetPost(ctx, slug)
		if err != nil {
			h.logger.Error("Error fetching post:", err)
			http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
			return
		}

		if post == nil {
			http.NotFound(w, r)
			return
		}

		// Handle different response types
		switch {
		case r.Header.Get("Accept") == "application/json":
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(post); err != nil {
				h.logger.Error("Error encoding post:", err)
				http.Error(w, "Error encoding response", http.StatusInternalServerError)
			}
		default:
			err = pages.BlogPost(post).Render(ctx, w)
			if err != nil {
				h.logger.Error("Error rendering post page:", err)
				http.Error(w, "Error rendering page", http.StatusInternalServerError)
			}
		}
	}
}

// CreatePost handles blog post creation
func (h *PostHandlers) CreatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var post models.Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			h.logger.Error("Error decoding post:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		if err := h.service.CreatePost(ctx, &post, []int64{}); err != nil {
			h.logger.Error("Error creating post:", err)
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(post)
	}
}

// UpdatePost handles blog post updates
func (h *PostHandlers) UpdatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut && r.Method != http.MethodPatch {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var post models.Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			h.logger.Error("Error decoding post:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get post ID from URL
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		post.ID = id

		ctx := r.Context()
		if err := h.service.UpdatePost(ctx, &post, []int64{}); err != nil {
			h.logger.Error("Error updating post:", err)
			http.Error(w, "Failed to update post", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}
}

// DeletePost handles blog post deletion
func (h *PostHandlers) DeletePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get post ID from URL
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		if err := h.service.DeletePost(ctx, id); err != nil {
			h.logger.Error("Error deleting post:", err)
			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
