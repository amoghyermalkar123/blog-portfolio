// internal/service/post_service.go
package service

import (
	"blog-portfolio/internal/models"
	"blog-portfolio/internal/repository"
	"blog-portfolio/internal/utils"
	"context"
	"strings"
	"time"
)

type PostService struct {
	repo *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

// CreatePost creates a new blog post
func (s *PostService) CreatePost(ctx context.Context, post *models.Post, tagIds []int64) error {
	// Generate slug if not provided
	if post.Slug == "" {
		post.Slug = generateSlug(post.Title)
	}

	// Set published time if post is published
	if post.Published && post.PublishedAt == nil {
		now := time.Now()
		post.PublishedAt = &now
	}

	return s.repo.CreatePost(ctx, post)
}

// GetPost retrieves a post by its slug
func (s *PostService) GetPost(ctx context.Context, slug string) (*models.Post, error) {
	post, err := s.repo.GetPost(ctx, slug)
	if err != nil {
		return nil, err
	}

	// Calculate reading time
	post.ReadingTime = utils.CalculateReadingTime(post.Content)

	return post, nil
}

// ListPosts returns a list of posts based on the filter
func (s *PostService) ListPosts(ctx context.Context, filter models.PostFilter) ([]*models.Post, error) {
	return s.repo.ListPosts(ctx, filter)
}

// UpdatePost updates an existing post
func (s *PostService) UpdatePost(ctx context.Context, post *models.Post, tagIds []int64) error {
	// Update published time if post is being published
	if post.Published && post.PublishedAt == nil {
		now := time.Now()
		post.PublishedAt = &now
	}

	return s.repo.UpdatePost(ctx, post)
}

// DeletePost deletes a post by ID
func (s *PostService) DeletePost(ctx context.Context, id int64) error {
	return s.repo.DeletePost(ctx, id)
}

// Helper function to generate URL-friendly slugs
func generateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)

	// Remove multiple consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// internal/service/post_service.go

// GetPostByID retrieves a post by its ID
func (s *PostService) GetPostByID(ctx context.Context, id int64) (*models.Post, error) {
	return s.repo.GetPostByID(ctx, id)
}

// ListTags returns all available tags
func (s *PostService) ListTags(ctx context.Context) ([]models.Tag, error) {
	return s.repo.ListTags(ctx)
}

// CreatePost creates a new post with tags

// UpdatePost updates an existing post and its tags

// generateSlug creates a URL-friendly version of the title
