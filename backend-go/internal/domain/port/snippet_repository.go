package port

import (
	"context"
	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type SnippetRepository interface {
	Save(ctx context.Context, snippet *model.Snippet) (*model.Snippet, error)
	FindByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (*model.Snippet, error)
	FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page int, size int) (model.Page[model.Snippet], error)
	DeleteByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (bool, error)
	ProjectScopedRepository
}
