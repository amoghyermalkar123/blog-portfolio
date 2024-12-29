// internal/models/post.go
package models

import (
	"time"
)

type Post struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Content     string     `json:"content"`
	Description string     `json:"description"`
	CoverImage  string     `json:"cover_image"`
	Published   bool       `json:"published"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	Tags        []Tag      `json:"tags,omitempty"`
}

// PostFilter represents filters for querying posts
type PostFilter struct {
	Tag       string
	Published *bool
	Limit     int
	Offset    int
}
