// internal/repository/tag_repository.go
package repository

import (
	"blog-portfolio/internal/models"
	"context"
	"database/sql"
	"strings"
)

type TagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) CreateTag(ctx context.Context, tag *models.Tag) error {
	query := `
        INSERT INTO tags (name, slug) 
        VALUES (?, ?)
        RETURNING id, created_at`

	// Generate slug from name
	slug := generateSlug(tag.Name)

	err := r.db.QueryRowContext(
		ctx,
		query,
		tag.Name,
		slug,
	).Scan(&tag.ID, &tag.CreatedAt)

	return err
}

func (r *TagRepository) UpdateTag(ctx context.Context, tag *models.Tag) error {
	query := `
        UPDATE tags 
        SET name = ?, slug = ?
        WHERE id = ?`

	// Generate new slug from updated name
	slug := generateSlug(tag.Name)

	result, err := r.db.ExecContext(ctx, query, tag.Name, slug, tag.ID)
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

func (r *TagRepository) DeleteTag(ctx context.Context, id int64) error {
	// Start a transaction since we need to delete from multiple tables
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// First delete from post_tags
	_, err = tx.ExecContext(ctx, "DELETE FROM post_tags WHERE tag_id = ?", id)
	if err != nil {
		return err
	}

	// Then delete the tag
	result, err := tx.ExecContext(ctx, "DELETE FROM tags WHERE id = ?", id)
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

	return tx.Commit()
}

func (r *TagRepository) ListTags(ctx context.Context) ([]models.Tag, error) {
	query := `
        SELECT t.id, t.name, t.slug, t.created_at, COUNT(pt.post_id) as post_count
        FROM tags t
        LEFT JOIN post_tags pt ON t.id = pt.tag_id
        GROUP BY t.id, t.name, t.slug, t.created_at
        ORDER BY t.name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		var postCount int
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Slug,
			&tag.CreatedAt,
			&postCount,
		)
		if err != nil {
			return nil, err
		}
		tag.PostCount = postCount
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *TagRepository) GetTagByID(ctx context.Context, id int64) (*models.Tag, error) {
	query := `
        SELECT id, name, slug, created_at
        FROM tags
        WHERE id = ?`

	var tag models.Tag
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tag.ID,
		&tag.Name,
		&tag.Slug,
		&tag.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

// Helper function to generate URL-friendly slugs
func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove any characters that aren't alphanumeric or hyphens
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
	return strings.Trim(slug, "-")
}
