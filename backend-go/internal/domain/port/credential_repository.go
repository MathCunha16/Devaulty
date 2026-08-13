package port

import (
	"context"
	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type CredentialRepository interface {
	Save(ctx context.Context, credential *model.Credential) (*model.Credential, error)
	FindByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (*model.Credential, error)
	FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page int, size int) (model.Page[model.Credential], error)
	DeleteByIDAndProjectID(ctx context.Context, ProjectID, id uuid.UUID) (bool, error)
	ProjectScopedRepository
}
