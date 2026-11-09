package port

import (
	"context"
	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type DrawingRepository interface {
	Save(ctx context.Context, drawing *model.Drawing) (*model.Drawing, error)
	FindByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (*model.Drawing, error)
	FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[model.Drawing], error)
	DeleteByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (bool, error)
	ProjectScopedRepository
}
