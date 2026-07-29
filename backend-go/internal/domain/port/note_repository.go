package port

import (
	"context"
	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type NoteRepository interface {
	ProjectScopedRepository
	Save(ctx context.Context, note *model.Note) (*model.Note, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Note, error)
	FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page int, size int) (model.Page[model.Note], error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
