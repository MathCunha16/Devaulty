package persistence

import (
	"context"
	"database/sql"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/domain/port"
	"devaulty-backend/internal/usecase"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AppSettingRepositoryAdapter struct {
	db *sqlx.DB
}

func NewAppSettingRepository(db *sqlx.DB) port.AppSettingRepository {
	return &AppSettingRepositoryAdapter{db: db}
}

func (r *AppSettingRepositoryAdapter) Save(ctx context.Context, setting *model.AppSetting) (*model.AppSetting, error) {
	query := `	INSERT INTO app_settings (key, value) 
				VALUES (:key, :value)
 				ON CONFLICT(key) DO UPDATE SET value = excluded.value`
	_, err := r.db.NamedExecContext(ctx, query, setting)
	if err != nil {
		return nil, fmt.Errorf("error trying to save app setting: %w", err)
	}

	return setting, nil
}

func (r *AppSettingRepositoryAdapter) FindByKey(ctx context.Context, key string) (*model.AppSetting, error) {
	query := `SELECT * FROM app_settings WHERE key = ?`
	var setting model.AppSetting
	err := r.db.GetContext(ctx, &setting, query, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // not found, returns nil
		}

		return nil, fmt.Errorf("error trying to find app setting: %w", err)
	}
	return &setting, nil
}

func (r *AppSettingRepositoryAdapter) ExistsByKey(ctx context.Context, key string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM app_settings WHERE key = ?)`
	var exits bool
	err := r.db.GetContext(ctx, &exits, query, key)
	if err != nil {
		return false, fmt.Errorf("error trying to check if app setting exists: %w", err)
	}
	return exits, nil
}

func (r *AppSettingRepositoryAdapter) SaveMasterPasswordSettings(ctx context.Context, hashValue, saltValue string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var count int
	queryCheck := `SELECT COUNT(*) FROM app_settings WHERE key IN ('master_password_hash', 'master_password_salt')`
	err = tx.GetContext(ctx, &count, queryCheck)
	if err != nil {
		return fmt.Errorf("error checking existing master password keys: %w", err)
	}
	if count > 0 {
		return usecase.ErrMasterPasswordAlreadyConfigured
	}

	queryInsert := `INSERT INTO app_settings (key, value) VALUES ('master_password_hash', ?), ('master_password_salt', ?)`
	_, err = tx.ExecContext(ctx, queryInsert, hashValue, saltValue)
	if err != nil {
		return fmt.Errorf("error inserting master password settings: %w", err)
	}

	return tx.Commit()
}
