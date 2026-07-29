package port

import (
	"context"
	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type ProblemRepository interface {
	ProjectScopedRepository
	Save(ctx context.Context, problem *model.Problem) (*model.Problem, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Problem, error)
	FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page int, size int) (model.Page[model.Problem], error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
