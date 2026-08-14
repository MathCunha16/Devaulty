package persistence

import (
	"context"
	"devaulty-backend/internal/domain/model"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func PaginateExec[T any](
	ctx context.Context,
	db *sqlx.DB,
	countQuery string,
	selectQuery string,
	page, size int,
	args ...any) (model.Page[T], error) {
	if page < 0 {
		page = 0
	}

	if size <= 0 {
		size = 10
	}

	offset := page * size
	var totalElements int64
	err := db.GetContext(ctx, &totalElements, countQuery, args...)
	if err != nil {
		return model.Page[T]{}, fmt.Errorf("error trying to count items: %w", err)
	}
	if totalElements == 0 {
		return model.NewPage([]T{}, page, size, 0), nil
	}

	queryWithPagination := fmt.Sprintf("%s LIMIT ? OFFSET ?", selectQuery)
	fullArgs := append(args, size, offset)
	var items []T
	err = db.SelectContext(ctx, &items, queryWithPagination, fullArgs...)
	if err != nil {
		return model.Page[T]{}, fmt.Errorf("error trying to find items: %w", err)
	}
	return model.NewPage(items, page, size, totalElements), nil
}
