package usecase

import (
	"context"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/domain/port"
	"devaulty-backend/internal/dto"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTagAlreadyExists = errors.New("tag already exists")
	ErrTagNotFound      = errors.New("tag not found")
)

type TagUseCase struct {
	tagRepo     port.TagRepository
	projectRepo port.ProjectRepository
}

func NewTagUseCase(tagRepo port.TagRepository, projectRepo port.ProjectRepository) *TagUseCase {
	return &TagUseCase{
		tagRepo:     tagRepo,
		projectRepo: projectRepo,
	}
}

func (uc *TagUseCase) Create(ctx context.Context, cmd dto.CreateTagCommand) (*model.Tag, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}

	sanitizedName := strings.TrimSpace(cmd.Name)
	alreadyExists, err := uc.tagRepo.ExistsByNameAndProjectID(ctx, sanitizedName, cmd.ProjectID)
	if err != nil {
		return nil, err
	}
	if alreadyExists {
		return nil, ErrTagAlreadyExists
	}

	tag := model.Tag{
		ID:        uuid.New(),
		ProjectID: cmd.ProjectID,
		Name:      sanitizedName,
		Color:     &cmd.Color,
		BaseEntity: model.BaseEntity{
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		},
	}
	return uc.tagRepo.Save(ctx, &tag)
}

func (uc *TagUseCase) GetByID(ctx context.Context, projectID, id uuid.UUID) (*model.Tag, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return nil, err
	}
	tag, err := uc.tagRepo.FindByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return nil, fmt.Errorf("error trying to find tag: %w", err)
	}
	if tag == nil {
		return nil, ErrTagNotFound
	}
	return tag, nil
}

func (uc *TagUseCase) GetAllByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.Tag, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return nil, err
	}
	return uc.tagRepo.FindAllByProjectID(ctx, projectID)
}

func (uc *TagUseCase) SearchByName(ctx context.Context, projectID uuid.UUID, name string) ([]model.Tag, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return nil, err
	}
	return uc.tagRepo.SearchByNameAndProjectID(ctx, name, projectID)
}

func (uc *TagUseCase) Update(ctx context.Context, cmd dto.UpdateTagCommand) (*model.Tag, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}

	tag, err := uc.tagRepo.FindByIDAndProjectID(ctx, cmd.ProjectID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error trying to find tag: %w", err)
	}
	if tag == nil {
		return nil, ErrTagNotFound
	}

	if cmd.Name != nil {
		sanitizedName := strings.TrimSpace(*cmd.Name)
		if !strings.EqualFold(tag.Name, sanitizedName) {
			alreadyExists, err := uc.tagRepo.ExistsByNameAndProjectID(ctx,
				sanitizedName, cmd.ProjectID)
			if err != nil {
				return nil, err
			}
			if alreadyExists {
				return nil, ErrTagAlreadyExists
			}
		}
		tag.Name = sanitizedName
	}
	if cmd.Color != nil {
		tag.Color = cmd.Color
	}
	now := time.Now()
	tag.UpdatedAt = &now
	return uc.tagRepo.Save(ctx, tag)
}

func (uc *TagUseCase) Delete(ctx context.Context, projectID, id uuid.UUID) error {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return err
	}
	deleted, err := uc.tagRepo.DeleteByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return fmt.Errorf("error deleting tag: %w", err)
	}
	if !deleted {
		return ErrTagNotFound
	}
	return nil
}
