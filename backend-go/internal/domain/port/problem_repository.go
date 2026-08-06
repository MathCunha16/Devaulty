package port

import (
	"context"
	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type ProblemSummary struct {
	ID        uuid.UUID             `json:"id" db:"id"`
	ProjectID uuid.UUID             `json:"projectId" db:"project_id"`
	Title     string                `json:"title" db:"title"`
	Status    model.ProblemStatus   `json:"status" db:"status"`
	Severity  model.ProblemSeverity `json:"severity" db:"severity"`
	model.BaseEntity
}

type ProblemRepository interface {
	Save(ctx context.Context, problem *model.Problem) (*model.Problem, error)
	FindByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (*model.Problem, error)
	FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page int, size int) (model.Page[ProblemSummary], error)
	DeleteByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (bool, error)
	ProjectScopedRepository
}
