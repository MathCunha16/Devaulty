package port

import (
	"context"
	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type SnippetRepository interface {
	Save(ctx context.Context, snippet *model.Snippet) (*model.Snippet, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Snippet, error)
	FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page int, size int) (model.Page[model.Snippet], error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	ProjectScopedRepository
}
