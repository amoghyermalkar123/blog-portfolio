// internal/handlers/admin_handlers.go
package handlers

import (
	"blog-portfolio/internal/logger"
	"blog-portfolio/internal/models"
	"blog-portfolio/internal/service"
	"blog-portfolio/web/pages/admin"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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

// ShowCreatePost handles displaying the new post form
func (h *AdminHandlers) ShowCreatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get all available tags
		tags, err := h.posts.ListTags(r.Context())
		if err != nil {
			h.logger.Error("Error fetching tags:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := admin.PostEditorData{
			IsNew: true,
			Tags:  tags,
		}

		err = admin.PostEditor(data).Render(r.Context(), w)
		if err != nil {
			h.logger.Error("Error rendering post editor:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// ShowEditPost handles displaying the edit post form
func (h *AdminHandlers) ShowEditPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get post ID from URL
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// Get post
		post, err := h.posts.GetPostByID(r.Context(), id)
		if err != nil {
			h.logger.Error("Error fetching post:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if post == nil {
			http.NotFound(w, r)
			return
		}

		// Get all available tags
		tags, err := h.posts.ListTags(r.Context())
		if err != nil {
			h.logger.Error("Error fetching tags:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := admin.PostEditorData{
			Post:  post,
			Tags:  tags,
			IsNew: false,
		}

		err = admin.PostEditor(data).Render(r.Context(), w)
		if err != nil {
			h.logger.Error("Error rendering post editor:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// HandleCreatePost processes the new post form submission
func (h *AdminHandlers) HandleCreatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse form
		if err := r.ParseForm(); err != nil {
			h.logger.Error("Error parsing form:", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Get action (draft or publish)
		action := r.FormValue("action")
		published := action == "publish"

		// Get selected tag IDs
		var tagIDs []int64
		for _, idStr := range r.Form["tags[]"] {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				continue
			}
			tagIDs = append(tagIDs, id)
		}

		// Create post
		post := &models.Post{
			Title:       r.FormValue("title"),
			Content:     r.FormValue("content"),
			Description: r.FormValue("description"),
			CoverImage:  r.FormValue("cover_image"),
			Published:   published,
		}

		if published {
			now := time.Now()
			post.PublishedAt = &now
		}

		// Save post
		err := h.posts.CreatePost(r.Context(), post, tagIDs)
		if err != nil {
			h.logger.Error("Error creating post:", err)
			// Re-render form with error
			tags, _ := h.posts.ListTags(r.Context())
			data := admin.PostEditorData{
				Post:  post,
				Tags:  tags,
				IsNew: true,
				Error: "Failed to create post: " + err.Error(),
			}
			admin.PostEditor(data).Render(r.Context(), w)
			return
		}

		// Redirect to the post list with success message
		http.Redirect(w, r, "/admin/posts?success=created", http.StatusSeeOther)
	}
}

// HandleUpdatePost processes the edit post form submission
func (h *AdminHandlers) HandleUpdatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get post ID from URL
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// Parse form
		if err := r.ParseForm(); err != nil {
			h.logger.Error("Error parsing form:", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Get existing post
		existingPost, err := h.posts.GetPostByID(r.Context(), id)
		if err != nil {
			h.logger.Error("Error fetching post:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if existingPost == nil {
			http.NotFound(w, r)
			return
		}

		// Get action (draft or publish)
		action := r.FormValue("action")
		published := action == "publish"

		// Get selected tag IDs
		var tagIDs []int64
		for _, idStr := range r.Form["tags[]"] {
			tagID, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				continue
			}
			tagIDs = append(tagIDs, tagID)
		}

		// Update post fields
		post := existingPost
		post.Title = r.FormValue("title")
		post.Content = r.FormValue("content")
		post.Description = r.FormValue("description")
		post.CoverImage = r.FormValue("cover_image")

		// Handle publication status change
		if published && !post.Published {
			// Post is being published for the first time
			now := time.Now()
			post.PublishedAt = &now
			post.Published = true
		} else if !published && post.Published {
			// Post is being unpublished
			post.Published = false
		}

		// Save updates
		err = h.posts.UpdatePost(r.Context(), post, tagIDs)
		if err != nil {
			h.logger.Error("Error updating post:", err)
			// Re-render form with error
			tags, _ := h.posts.ListTags(r.Context())
			data := admin.PostEditorData{
				Post:  post,
				Tags:  tags,
				IsNew: false,
				Error: "Failed to update post: " + err.Error(),
			}
			admin.PostEditor(data).Render(r.Context(), w)
			return
		}

		// Redirect to the post list with success message
		http.Redirect(w, r, "/admin/posts?success=updated", http.StatusSeeOther)
	}
}
