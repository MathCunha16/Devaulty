package port

import (
	"context"
	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type TagRepository interface {
	Save(ctx context.Context, tag *model.Tag) (*model.Tag, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Tag, error)
	FindAllByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.Tag, error)
	SearchByNameAndProjectID(ctx context.Context, name string, projectID uuid.UUID) ([]model.Tag, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	ExistsByNameAndProjectID(ctx context.Context, name string, projectID uuid.UUID) (bool, error)
	ExistsByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (bool, error)
}
