package port

import (
	"context"
	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type ProjectScopedRepository interface {
	GetSupportedType() model.ItemType
	ExistsByIdAndProjectID(ctx context.Context, id uuid.UUID, projectId uuid.UUID) (bool, error)
	FindExistingIdsByProjectID(ctx context.Context, ids []uuid.UUID, projectId uuid.UUID) ([]uuid.UUID, error)
}
