// internal/models/post.go
package models

import (
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
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
	ReadingTime int        `json:"reading_time"`
}

// PostFilter represents filters for querying posts
type PostFilter struct {
	Tag       string
	Published *bool
	Limit     int
	Offset    int
}

func (p *Post) ParsedContent() string {
	// Create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	post := parser.NewWithExtensions(extensions)

	// Parse the markdown
	doc := post.Parse([]byte(p.Content))

	// Create HTML renderer with options
	opts := html.RendererOptions{
		Flags: html.CommonFlags | html.HrefTargetBlank,
		Title: p.Title,
	}
	renderer := html.NewRenderer(opts)

	// Render to HTML
	return string(markdown.Render(doc, renderer))
}
