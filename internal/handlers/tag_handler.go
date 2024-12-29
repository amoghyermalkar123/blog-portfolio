// internal/handlers/tag_handlers.go
package handlers

import (
	"blog-portfolio/internal/logger"
	"blog-portfolio/internal/models"
	"blog-portfolio/internal/service"
	"blog-portfolio/web/pages/admin"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type TagHandlers struct {
	logger *logger.Logger
	tags   *service.TagService
}

func NewTagHandlers(logger *logger.Logger, tagService *service.TagService) *TagHandlers {
	return &TagHandlers{
		logger: logger,
		tags:   tagService,
	}
}

// ShowTags displays the tag management page
func (h *TagHandlers) ShowTags() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tags, err := h.tags.ListTags(ctx)
		if err != nil {
			h.logger.Error("Error fetching tags:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := admin.TagListData{
			Tags: tags,
		}

		err = admin.Tags(data).Render(ctx, w)
		if err != nil {
			h.logger.Error("Error rendering tags page:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// CreateTag handles tag creation
func (h *TagHandlers) CreateTag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateTagRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		tag, err := h.tags.CreateTag(r.Context(), &req)
		if err != nil {
			h.logger.Error("Error creating tag:", err)
			http.Error(w, "Failed to create tag", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tag)
	}
}

// UpdateTag handles tag updates
func (h *TagHandlers) UpdateTag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid tag ID", http.StatusBadRequest)
			return
		}

		var req models.UpdateTagRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = h.tags.UpdateTag(r.Context(), id, &req)
		if err != nil {
			h.logger.Error("Error updating tag:", err)
			http.Error(w, "Failed to update tag", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// DeleteTag handles tag deletion
func (h *TagHandlers) DeleteTag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid tag ID", http.StatusBadRequest)
			return
		}

		err = h.tags.DeleteTag(r.Context(), id)
		if err != nil {
			h.logger.Error("Error deleting tag:", err)
			http.Error(w, "Failed to delete tag", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
