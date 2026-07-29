package port

import (
	"context"
	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type CredentialRepository interface {
	ProjectScopedRepository
	Save(ctx context.Context, credential *model.Credential) (*model.Credential, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Credential, error)
	FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page int, size int) (model.Page[model.Credential], error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
