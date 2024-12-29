// internal/models/tag.go
package models

import "time"

type Tag struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	PostCount int       `json:"post_count"` // This will be populated when listing tags
}

// Request/Response structures
type CreateTagRequest struct {
	Name string `json:"name"`
}

type UpdateTagRequest struct {
	Name string `json:"name"`
}

type TagListResponse struct {
	Tags []Tag `json:"tags"`
}
