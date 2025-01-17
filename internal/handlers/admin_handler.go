// internal/handlers/admin_handlers.go
package handlers

import (
	"blog-portfolio/internal/logger"
	"blog-portfolio/internal/models"
	"blog-portfolio/internal/service"
	"blog-portfolio/internal/utils"
	"blog-portfolio/web/pages/admin"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type AdminHandlers struct {
	logger *logger.Logger
	posts  *service.PostService
	tags   *service.TagService
}

func NewAdminHandlers(logger *logger.Logger, postService *service.PostService, tagService *service.TagService) *AdminHandlers {
	return &AdminHandlers{
		logger: logger,
		posts:  postService,
		tags:   tagService,
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
		published := action == "publish" // This is correct but let's add logging

		h.logger.Info("Post action:", action, "Published:", published) // Add logging

		// Get selected tag IDs
		var tagIDs []int64
		for _, idStr := range r.Form["tags[]"] {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				continue
			}
			tagIDs = append(tagIDs, id)
		}

		// Create post with proper publishing status
		post := &models.Post{
			Title:       r.FormValue("title"),
			Content:     r.FormValue("content"),
			Description: r.FormValue("description"),
			CoverImage:  r.FormValue("cover_image"),
			Published:   published,
		}

		// Set published date if being published
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

func (h *AdminHandlers) HandleDeletePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get post ID from URL
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			h.logger.Error("Invalid post ID:", err)
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// Delete post
		err = h.posts.DeletePost(r.Context(), id)
		if err != nil {
			h.logger.Error("Error deleting post:", err)
			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
			return
		}

		// Return 200 OK - HTMX will handle removing the element from the DOM
		w.WriteHeader(http.StatusOK)
	}
}

// Inside internal/handlers/admin_handler.go

// Inside internal/handlers/admin_handler.go

// Inside internal/handlers/admin_handler.go

func (h *AdminHandlers) HandlePreview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			h.logger.Error("Error parsing form:", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Debug log the form values
		h.logger.Info("Preview form values:",
			"title:", r.FormValue("title"),
			"description:", r.FormValue("description"),
			"content:", "content length: "+strconv.Itoa(len(r.FormValue("content"))))

		// Create a temporary post from form data
		now := time.Now()
		post := &models.Post{
			Title:       r.FormValue("title"),
			Content:     r.FormValue("content"),
			Description: r.FormValue("description"),
			CoverImage:  r.FormValue("cover_image"),
			CreatedAt:   now,
			UpdatedAt:   now,
			PublishedAt: &now, // Add this for preview
		}

		// Parse tags if present
		var tags []models.Tag
		for _, idStr := range r.Form["tags[]"] {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				continue
			}
			tag, err := h.tags.GetTagByID(r.Context(), id)
			if err != nil {
				continue
			}
			if tag != nil {
				tags = append(tags, *tag)
			}
		}
		post.Tags = tags

		// Calculate reading time
		post.ReadingTime = utils.CalculateReadingTime(post.Content)

		// Debug log the post object
		h.logger.Info("Preview post object:",
			"title:", post.Title,
			"hasDescription:", post.Description != "",
			"contentLength:", len(post.Content))

		// Render preview page
		err := admin.Preview(admin.PreviewData{
			Post: post,
		}).Render(r.Context(), w)
		if err != nil {
			h.logger.Error("Error rendering preview:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
