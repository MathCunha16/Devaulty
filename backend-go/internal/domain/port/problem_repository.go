package port

import (
	"context"
	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type ProblemRepository interface {
	Save(ctx context.Context, problem *model.Problem) (*model.Problem, error)
	FindByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (*model.Problem, error)
	FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page int, size int) (model.Page[model.Problem], error)
	DeleteByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (bool, error)
	ProjectScopedRepository
}
