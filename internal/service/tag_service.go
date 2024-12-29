// internal/service/tag_service.go
package service

import (
	"blog-portfolio/internal/models"
	"blog-portfolio/internal/repository"
	"context"
)

type TagService struct {
	repo *repository.TagRepository
}

func NewTagService(repo *repository.TagRepository) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) CreateTag(ctx context.Context, tagRequest *models.CreateTagRequest) (*models.Tag, error) {
	tag := &models.Tag{
		Name: tagRequest.Name,
	}

	err := s.repo.CreateTag(ctx, tag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *TagService) UpdateTag(ctx context.Context, id int64, tagRequest *models.UpdateTagRequest) error {
	tag := &models.Tag{
		ID:   id,
		Name: tagRequest.Name,
	}

	return s.repo.UpdateTag(ctx, tag)
}

func (s *TagService) DeleteTag(ctx context.Context, id int64) error {
	return s.repo.DeleteTag(ctx, id)
}

func (s *TagService) ListTags(ctx context.Context) ([]models.Tag, error) {
	return s.repo.ListTags(ctx)
}

func (s *TagService) GetTagByID(ctx context.Context, id int64) (*models.Tag, error) {
	return s.repo.GetTagByID(ctx, id)
}
