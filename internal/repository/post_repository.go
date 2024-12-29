// internal/repository/post_repository.go
package repository

import (
	"blog-portfolio/internal/models"
	"context"
	"database/sql"
	"errors"
	"strings"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

// CreatePost creates a new blog post
func (r *PostRepository) CreatePost(ctx context.Context, post *models.Post) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert post
	query := `
        INSERT INTO posts (title, slug, content, description, cover_image, published, published_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        RETURNING id, created_at, updated_at`

	var publishedAt sql.NullTime
	if post.PublishedAt != nil {
		publishedAt.Time = *post.PublishedAt
		publishedAt.Valid = true
	}

	err = tx.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Slug,
		post.Content,
		post.Description,
		post.CoverImage,
		post.Published,
		publishedAt,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}

	// Insert tags if any
	if len(post.Tags) > 0 {
		for _, tag := range post.Tags {
			// Insert tag if it doesn't exist
			var tagID int64
			err = tx.QueryRowContext(
				ctx,
				`INSERT INTO tags (name, slug) 
                 VALUES (?, ?) 
                 ON CONFLICT(slug) DO UPDATE SET name=excluded.name 
                 RETURNING id`,
				tag.Name,
				tag.Slug,
			).Scan(&tagID)
			if err != nil {
				return err
			}

			// Create post-tag association
			_, err = tx.ExecContext(
				ctx,
				"INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)",
				post.ID,
				tagID,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// GetPost retrieves a post by its slug
func (r *PostRepository) GetPost(ctx context.Context, slug string) (*models.Post, error) {
	post := &models.Post{}
	query := `
        SELECT 
            p.id, p.title, p.slug, p.content, p.description, 
            p.cover_image, p.published, p.created_at, p.updated_at, 
            p.published_at 
        FROM posts p 
        WHERE p.slug = ?`

	var publishedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&post.ID,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.Description,
		&post.CoverImage,
		&post.Published,
		&post.CreatedAt,
		&post.UpdatedAt,
		&publishedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if publishedAt.Valid {
		post.PublishedAt = &publishedAt.Time
	}

	// Get tags
	post.Tags, err = r.getPostTags(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// ListPosts returns a list of posts based on the filter
// internal/repository/post_repository.go

func (r *PostRepository) ListPosts(ctx context.Context, filter models.PostFilter) ([]*models.Post, error) {
	query := strings.Builder{}
	query.WriteString(`
        SELECT DISTINCT
            p.id, p.title, p.slug, p.content, p.description, 
            p.cover_image, p.published, p.created_at, 
            p.updated_at, p.published_at
        FROM posts p
    `)

	args := []interface{}{}
	where := []string{}

	if filter.Tag != "" {
		query.WriteString(` 
            LEFT JOIN post_tags pt ON p.id = pt.post_id 
            LEFT JOIN tags t ON pt.tag_id = t.id
        `)
		where = append(where, "t.slug = ?")
		args = append(args, filter.Tag)
	}

	// Important: Add published filter condition
	if filter.Published != nil {
		where = append(where, "p.published = ?")
		args = append(args, *filter.Published)
	}

	if len(where) > 0 {
		query.WriteString(" WHERE " + strings.Join(where, " AND "))
	}

	// Order by published date for published posts, creation date for drafts
	query.WriteString(" ORDER BY CASE WHEN p.published = 1 THEN p.published_at ELSE p.created_at END DESC")

	if filter.Limit > 0 {
		query.WriteString(" LIMIT ?")
		args = append(args, filter.Limit)
	}
	if filter.Offset > 0 {
		query.WriteString(" OFFSET ?")
		args = append(args, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		var publishedAt sql.NullTime
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Slug,
			&post.Content,
			&post.Description,
			&post.CoverImage,
			&post.Published,
			&post.CreatedAt,
			&post.UpdatedAt,
			&publishedAt,
		)
		if err != nil {
			return nil, err
		}

		if publishedAt.Valid {
			post.PublishedAt = &publishedAt.Time
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// UpdatePost updates an existing post
func (r *PostRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update post
	query := `
        UPDATE posts 
        SET title = ?, content = ?, description = ?, 
            cover_image = ?, published = ?, published_at = ?,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = ?`

	var publishedAt sql.NullTime
	if post.PublishedAt != nil {
		publishedAt.Time = *post.PublishedAt
		publishedAt.Valid = true
	}

	result, err := tx.ExecContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.Description,
		post.CoverImage,
		post.Published,
		publishedAt,
		post.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	// Update tags
	_, err = tx.ExecContext(ctx, "DELETE FROM post_tags WHERE post_id = ?", post.ID)
	if err != nil {
		return err
	}

	for _, tag := range post.Tags {
		var tagID int64
		err = tx.QueryRowContext(
			ctx,
			`INSERT INTO tags (name, slug) 
             VALUES (?, ?) 
             ON CONFLICT(slug) DO UPDATE SET name=excluded.name 
             RETURNING id`,
			tag.Name,
			tag.Slug,
		).Scan(&tagID)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)",
			post.ID,
			tagID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// DeletePost deletes a post by ID
func (r *PostRepository) DeletePost(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// getPostTags retrieves all tags for a given post
func (r *PostRepository) getPostTags(ctx context.Context, postID int64) ([]models.Tag, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT t.id, t.name, t.slug, t.created_at
        FROM tags t
        JOIN post_tags pt ON t.id = pt.tag_id
        WHERE pt.post_id = ?
        ORDER BY t.name`,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Slug, &tag.CreatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// internal/repository/post_repository.go

// BeginTx starts a new transaction
func (r *PostRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

// GetPostByID retrieves a post by its ID
func (r *PostRepository) GetPostByID(ctx context.Context, id int64) (*models.Post, error) {
	post := &models.Post{}
	query := `
        SELECT id, title, slug, content, description, cover_image, 
               published, created_at, updated_at, published_at
        FROM posts
        WHERE id = ?`

	var publishedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.Description,
		&post.CoverImage,
		&post.Published,
		&post.CreatedAt,
		&post.UpdatedAt,
		&publishedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if publishedAt.Valid {
		post.PublishedAt = &publishedAt.Time
	}

	// Get tags
	post.Tags, err = r.getPostTags(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// ListTags returns all available tags
func (r *PostRepository) ListTags(ctx context.Context) ([]models.Tag, error) {
	query := `
        SELECT id, name, slug, created_at
        FROM tags
        ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Slug, &tag.CreatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// CreatePostTx creates a new post within a transaction
func (r *PostRepository) CreatePostTx(ctx context.Context, tx *sql.Tx, post *models.Post) error {
	query := `
        INSERT INTO posts (title, slug, content, description, cover_image, published, published_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        RETURNING id, created_at, updated_at`

	var publishedAt sql.NullTime
	if post.PublishedAt != nil {
		publishedAt.Time = *post.PublishedAt
		publishedAt.Valid = true
	}

	err := tx.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Slug,
		post.Content,
		post.Description,
		post.CoverImage,
		post.Published,
		publishedAt,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	return err
}

// UpdatePostTx updates an existing post within a transaction
func (r *PostRepository) UpdatePostTx(ctx context.Context, tx *sql.Tx, post *models.Post) error {
	query := `
        UPDATE posts 
        SET title = ?, content = ?, description = ?, 
            cover_image = ?, published = ?, published_at = ?,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = ?`

	var publishedAt sql.NullTime
	if post.PublishedAt != nil {
		publishedAt.Time = *post.PublishedAt
		publishedAt.Valid = true
	}

	result, err := tx.ExecContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.Description,
		post.CoverImage,
		post.Published,
		publishedAt,
		post.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// SetPostTagsTx updates a post's tags within a transaction
func (r *PostRepository) SetPostTagsTx(ctx context.Context, tx *sql.Tx, postID int64, tagIDs []int64) error {
	// First, remove all existing tag associations
	_, err := tx.ExecContext(ctx, "DELETE FROM post_tags WHERE post_id = ?", postID)
	if err != nil {
		return err
	}

	// Then add new tag associations
	if len(tagIDs) > 0 {
		query := "INSERT INTO post_tags (post_id, tag_id) VALUES "
		var values []interface{}
		placeholders := make([]string, 0, len(tagIDs))

		for _, tagID := range tagIDs {
			placeholders = append(placeholders, "(?, ?)")
			values = append(values, postID, tagID)
		}

		query += strings.Join(placeholders, ", ")
		_, err = tx.ExecContext(ctx, query, values...)
		if err != nil {
			return err
		}
	}

	return nil
}
