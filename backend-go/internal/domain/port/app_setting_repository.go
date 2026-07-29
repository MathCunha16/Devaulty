package port

import (
	"context"
	"devaulty-backend/internal/domain/model"
)

type AppSettingRepository interface {
	Save(ctx context.Context, setting *model.AppSetting) (*model.AppSetting, error)
	FindByKey(ctx context.Context, key string) (*model.AppSetting, error)
	ExistsByKey(ctx context.Context, key string) (bool, error)
}
