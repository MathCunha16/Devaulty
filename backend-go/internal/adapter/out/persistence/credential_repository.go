package persistence

import (
	"context"
	"database/sql"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/domain/port"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CredentialRepositoryAdapter struct {
	db *sqlx.DB
}

func NewCredentialRepository(db *sqlx.DB) port.CredentialRepository {
	return &CredentialRepositoryAdapter{db: db}
}

func (r *CredentialRepositoryAdapter) Save(ctx context.Context, credential *model.Credential) (*model.Credential, error) {
	query := `	INSERT INTO credentials (id, project_id, title, secret_type, payload_encrypted,
	            	encryption_iv, encryption_auth_tag, notes,
	                related_url, created_at, updated_at) 
				VALUES (:id , :project_id, :title,
	        		:secret_type, :payload_encrypted, :encryption_iv, :encryption_auth_tag,
					:notes, :related_url, :created_at, :updated_at)
				ON CONFLICT(id) DO UPDATE SET title = excluded.title, secret_type = excluded.secret_type,
					payload_encrypted = excluded.payload_encrypted, encryption_iv = excluded.encryption_iv,
					encryption_auth_tag = excluded.encryption_auth_tag, notes = excluded.notes, 
					related_url = excluded.related_url, updated_at = excluded.updated_at`
	_, err := r.db.NamedExecContext(ctx, query, credential)
	if err != nil {
		return nil, fmt.Errorf("error trying to save credential: %w", err)
	}
	return credential, nil
}

func (r *CredentialRepositoryAdapter) FindByID(ctx context.Context, id uuid.UUID) (*model.Credential, error) {
	query := `SELECT * FROM credentials WHERE id = ?`
	var credential model.Credential
	err := r.db.GetContext(ctx, &credential, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error trying to find credential: %w", err)
	}
	return &credential, nil
}

func (r *CredentialRepositoryAdapter) FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page int, size int) (model.Page[model.Credential], error) {
	countQuery := `SELECT COUNT(*) FROM credentials WHERE project_id = ?`
	selectQuery := `SELECT * FROM credentials WHERE project_id = ? ORDER BY created_at DESC`
	return PaginateExec[model.Credential](ctx, r.db, countQuery, selectQuery, page, size, projectID)
}

func (r *CredentialRepositoryAdapter) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM credentials WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error trying to delete credential: %w", err)
	}
	return nil
}

func (r *CredentialRepositoryAdapter) ExistsByIDAndProjectID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (bool, error) {
	return ExistsByIDAndProjectIDExec(ctx, r.db, "credentials", id, projectID)
}

func (r *CredentialRepositoryAdapter) FindExistingIDsByProjectID(ctx context.Context, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	return FindExistingIDsByProjectIDExec(ctx, r.db, "credentials", ids, projectID)
}
