// internal/handlers/admin_handlers.go
package handlers

import (
	"blog-portfolio/internal/logger"
	"blog-portfolio/internal/models"
	"blog-portfolio/internal/service"
	"blog-portfolio/web/pages/admin"
	"net/http"
	"strconv"
)

type AdminHandlers struct {
	logger *logger.Logger
	posts  *service.PostService
}

func NewAdminHandlers(logger *logger.Logger, postService *service.PostService) *AdminHandlers {
	return &AdminHandlers{
		logger: logger,
		posts:  postService,
	}
}

// ShowDashboard handles the admin dashboard page
func (h *AdminHandlers) ShowDashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get post statistics
		filter := models.PostFilter{
			Limit: 5, // Only get recent 5 posts
		}
		recentPosts, err := h.posts.ListPosts(ctx, filter)
		if err != nil {
			h.logger.Error("Error fetching recent posts:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Count total posts
		totalPosts := len(recentPosts)
		publishedPosts := 0
		draftPosts := 0
		for _, post := range recentPosts {
			if post.Published {
				publishedPosts++
			} else {
				draftPosts++
			}
		}

		data := admin.DashboardData{
			PostCount:      totalPosts,
			PublishedCount: publishedPosts,
			DraftCount:     draftPosts,
			RecentPosts:    recentPosts,
		}

		err = admin.Dashboard(data).Render(ctx, w)
		if err != nil {
			h.logger.Error("Error rendering dashboard:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// ShowPosts handles the post listing page
func (h *AdminHandlers) ShowPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get page number from query params
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}

		// Set up pagination
		limit := 10
		offset := (page - 1) * limit

		// Get posts
		filter := models.PostFilter{
			Limit:  limit,
			Offset: offset,
		}
		posts, err := h.posts.ListPosts(ctx, filter)
		if err != nil {
			h.logger.Error("Error fetching posts:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// TODO: Get total count for pagination
		totalPages := 1 // This should be calculated based on total posts

		data := admin.PostListData{
			Posts:       posts,
			CurrentPage: page,
			TotalPages:  totalPages,
		}

		err = admin.Posts(data).Render(ctx, w)
		if err != nil {
			h.logger.Error("Error rendering posts page:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
