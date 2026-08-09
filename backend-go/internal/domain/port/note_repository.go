package port

import (
	"context"
	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type NoteRepository interface {
	Save(ctx context.Context, note *model.Note) (*model.Note, error)
	FindByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (*model.Note, error)
	FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page int, size int) (model.Page[model.Note], error)
	DeleteByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (bool, error)
	ProjectScopedRepository
}
