package usecase

import (
	"context"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/domain/port"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ProjectUseCase struct {
	projectRepo port.ProjectRepository
}

type CreateProjectCommand struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
	Icon        *string `json:"icon,omitempty"`
	Color       *string `json:"color,omitempty"`
}

type UpdateProjectCommand struct {
	ID          uuid.UUID // Is passed as a path variable in the controller
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	Icon        *string   `json:"icon,omitempty"`
	Color       *string   `json:"color,omitempty"`
}

func NewProjectUseCase(projectRepo port.ProjectRepository) *ProjectUseCase {
	return &ProjectUseCase{projectRepo: projectRepo}
}

func (uc *ProjectUseCase) Create(ctx context.Context, cmd CreateProjectCommand) (*model.Project, error) {
	project := model.Project{
		ID:          uuid.New(),
		Name:        cmd.Name,
		Description: cmd.Description,
		Icon:        cmd.Icon,
		Color:       cmd.Color,
		Archived:    false,
		BaseEntity: model.BaseEntity{
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		},
	}
	return uc.projectRepo.Save(ctx, &project)
}

func (uc *ProjectUseCase) GetByID(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	return uc.projectRepo.FindByID(ctx, id)
}

func (uc *ProjectUseCase) GetAll(ctx context.Context, page, size int) (model.Page[model.Project], error) {
	return uc.projectRepo.FindAll(ctx, page, size)
}

func (uc *ProjectUseCase) Update(ctx context.Context, cmd UpdateProjectCommand) (*model.Project, error) {
	project, err := uc.projectRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error trying to find project: %w", err)
	}
	if project == nil {
		return nil, errors.New("project not found")
	}

	if cmd.Name != nil {
		project.Name = *cmd.Name
	}
	if cmd.Description != nil {
		project.Description = cmd.Description
	}
	if cmd.Icon != nil {
		project.Icon = cmd.Icon
	}
	if cmd.Color != nil {
		project.Color = cmd.Color
	}
	now := time.Now()
	project.UpdatedAt = &now
	return uc.projectRepo.Save(ctx, project)
}

func (uc *ProjectUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.projectRepo.DeleteByID(ctx, id)
}

func (uc *ProjectUseCase) Archive(ctx context.Context, id uuid.UUID) error {
	project, err := uc.projectRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error trying to find project: %w", err)
	}
	if project == nil {
		return errors.New("project not found")
	}
	if project.Archived {
		return errors.New("project already archived")
	}

	project.Archived = true
	now := time.Now()
	project.UpdatedAt = &now
	_, err = uc.projectRepo.Save(ctx, project)
	if err != nil {
		return fmt.Errorf("error trying to archive project: %w", err)
	}
	return nil
}

func (uc *ProjectUseCase) Unarchive(ctx context.Context, id uuid.UUID) error {
	project, err := uc.projectRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error trying to find project: %w", err)
	}
	if project == nil {
		return errors.New("project not found")
	}
	if !project.Archived {
		return errors.New("project already unarchived")
	}

	project.Archived = false
	now := time.Now()
	project.UpdatedAt = &now
	_, err = uc.projectRepo.Save(ctx, project)
	if err != nil {
		return fmt.Errorf("error trying to unarchive project: %w", err)
	}
	return nil
}
