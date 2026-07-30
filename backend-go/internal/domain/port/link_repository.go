package port

import (
	"context"
	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type LinkRepository interface {
	Save(ctx context.Context, link *model.Link) (*model.Link, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Link, error)
	FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page int, size int) (model.Page[model.Link], error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	ProjectScopedRepository
}
