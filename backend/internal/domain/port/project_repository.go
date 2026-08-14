package port

import (
	"context"
	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type ProjectRepository interface {
	Save(ctx context.Context, project *model.Project) (*model.Project, error)
	FindAll(ctx context.Context, page int, size int) (model.Page[model.Project], error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Project, error)
	DeleteByID(ctx context.Context, id uuid.UUID) (bool, error)
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
}
