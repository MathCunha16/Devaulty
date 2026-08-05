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

var (
	ErrProblemNotFound = errors.New("problem not found")
)

type ProblemUseCase struct {
	problemRepo port.ProblemRepository
	projectRepo port.ProjectRepository
}

type CreateProblemCommand struct {
	ProjectID        uuid.UUID
	Title            string                `json:"title" binding:"required,min=2,max=255"`
	ErrorDescription string                `json:"errorDescription" binding:"required,min=2,max=255"`
	Solution         *string               `json:"solution,omitempty" binding:"omitempty,min=2,max=255"`
	Status           model.ProblemStatus   `json:"status" binding:"required"`
	Severity         model.ProblemSeverity `json:"severity" binding:"required"`
}

type UpdateProblemCommand struct {
	ProjectID        uuid.UUID
	ID               uuid.UUID
	Title            *string                `json:"title,omitempty" binding:"omitempty,min=2,max=255"`
	ErrorDescription *string                `json:"errorDescription,omitempty" binding:"omitempty,min=2,max=255"`
	Solution         *string                `json:"solution,omitempty" binding:"omitempty,min=2,max=255"`
	Severity         *model.ProblemSeverity `json:"severity,omitempty" binding:"omitempty"`
}

type UpdateProblemStatusCommand struct {
	ProjectID uuid.UUID
	ID        uuid.UUID
	Status    model.ProblemStatus `json:"status" binding:"required"`
}

func NewProblemUseCase(problemRepo port.ProblemRepository, projectRepo port.ProjectRepository) *ProblemUseCase {
	return &ProblemUseCase{
		problemRepo: problemRepo,
		projectRepo: projectRepo,
	}
}

func (uc *ProblemUseCase) Create(ctx context.Context, cmd CreateProblemCommand) (*model.Problem, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}

	problem := model.Problem{
		ID:               uuid.New(),
		ProjectID:        cmd.ProjectID,
		Title:            cmd.Title,
		ErrorDescription: cmd.ErrorDescription,
		Solution:         cmd.Solution,
		Status:           cmd.Status,
		Severity:         cmd.Severity,
		BaseEntity: model.BaseEntity{
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		},
	}

	return uc.problemRepo.Save(ctx, &problem)
}

func (uc *ProblemUseCase) GetByID(ctx context.Context, projectID, id uuid.UUID) (*model.Problem, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return nil, err
	}
	problem, err := uc.problemRepo.FindByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return nil, fmt.Errorf("error trying to find problem: %w", err)
	}
	if problem == nil {
		return nil, ErrProblemNotFound
	}
	return problem, nil
}

func (uc *ProblemUseCase) GetAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[model.Problem], error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return model.Page[model.Problem]{}, err
	}
	return uc.problemRepo.FindAllByProjectID(ctx, projectID, page, size)
}

func (uc *ProblemUseCase) Update(ctx context.Context, cmd UpdateProblemCommand) (*model.Problem, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}
	problem, err := uc.problemRepo.FindByIDAndProjectID(ctx, cmd.ProjectID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error trying to find problem: %w", err)
	}
	if problem == nil {
		return nil, ErrProblemNotFound
	}

	if cmd.Title != nil {
		problem.Title = *cmd.Title
	}
	if cmd.ErrorDescription != nil {
		problem.ErrorDescription = *cmd.ErrorDescription
	}
	if cmd.Solution != nil {
		problem.Solution = cmd.Solution
	}
	if cmd.Severity != nil {
		problem.Severity = *cmd.Severity
	}
	now := time.Now()
	problem.UpdatedAt = &now
	return uc.problemRepo.Save(ctx, problem)
}

func (uc *ProblemUseCase) UpdateStatus(ctx context.Context, cmd UpdateProblemStatusCommand) (*model.Problem, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}
	problem, err := uc.problemRepo.FindByIDAndProjectID(ctx, cmd.ProjectID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error trying to find problem: %w", err)
	}
	if problem == nil {
		return nil, ErrProblemNotFound
	}

	problem.Status = cmd.Status
	now := time.Now()
	problem.UpdatedAt = &now
	return uc.problemRepo.Save(ctx, problem)
}

func (uc *ProblemUseCase) Delete(ctx context.Context, projectID, id uuid.UUID) error {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return err
	}

	deleted, err := uc.problemRepo.DeleteByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return fmt.Errorf("error deleting problem: %w", err)
	}
	if !deleted {
		return ErrProblemNotFound
	}
	return nil
}
