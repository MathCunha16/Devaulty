package port

import (
	"context"
	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type BoardRepository interface {
	Save(ctx context.Context, board *model.Board) (*model.Board, error)
	FindByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (*model.Board, error)
	FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[model.Board], error)
	FindDefaultByProjectID(ctx context.Context, projectID uuid.UUID) (*model.Board, error)
	DeleteByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (bool, error)
	UnsetAllDefaultsByProjectID(ctx context.Context, projectID uuid.UUID) (bool, error)
	ProjectScopedRepository
}
