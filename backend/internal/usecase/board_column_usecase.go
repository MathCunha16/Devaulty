package usecase

import (
	"context"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/domain/port"
	"devaulty-backend/internal/dto"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrBoardColumnNotFound                     = errors.New("board column not found")
	ErrNotPossibleToGetBoardColumnNextPosition = errors.New("not possible to get board column next position")
)

type BoardColumnUseCase struct {
	boardColumnRepo port.BoardColumnRepository
	boardRepo       port.BoardRepository
	projectRepo     port.ProjectRepository
}

func NewBoardColumnUseCase(boardColumnRepo port.BoardColumnRepository, boardRepo port.BoardRepository, projectRepo port.ProjectRepository) *BoardColumnUseCase {
	return &BoardColumnUseCase{
		boardColumnRepo: boardColumnRepo,
		boardRepo:       boardRepo,
		projectRepo:     projectRepo,
	}
}

func (uc *BoardColumnUseCase) Create(ctx context.Context, cmd dto.CreateBoardColumnCommand) (*dto.BoardColumnView, error) {
	err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID)
	if err != nil {
		return nil, err
	}
	err = ensureBoardExists(ctx, uc.boardRepo, cmd.ProjectID, cmd.BoardID)
	if err != nil {
		return nil, err
	}

	nextPosition, err := uc.boardColumnRepo.GetNextPosition(ctx, cmd.BoardID)
	if err != nil {
		return nil, ErrNotPossibleToGetBoardColumnNextPosition
	}

	bc := model.BoardColumn{
		ID:       uuid.New(),
		BoardID:  cmd.BoardID,
		Name:     cmd.Name,
		Position: nextPosition,
		WipLimit: cmd.WipLimit,
		BaseEntity: model.BaseEntity{
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		},
	}

	saved, err := uc.boardColumnRepo.Save(ctx, &bc)
	if err != nil {
		return nil, err
	}

	return mapToBoardColumnView(saved, cmd.ProjectID), nil
}

func (uc *BoardColumnUseCase) GetByIDAndBoardID(ctx context.Context, projectID, boardID, id uuid.UUID) (*dto.BoardColumnView, error) {
	err := ensureProjectExists(ctx, uc.projectRepo, projectID)
	if err != nil {
		return nil, err
	}
	err = ensureBoardExists(ctx, uc.boardRepo, projectID, boardID)
	if err != nil {
		return nil, err
	}

	bc, err := uc.boardColumnRepo.FindByIDAndBoardID(ctx, boardID, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching board column: %w", err)
	}
	if bc == nil {
		return nil, ErrBoardColumnNotFound
	}
	return mapToBoardColumnView(bc, projectID), nil
}

func (uc *BoardColumnUseCase) GetAllByBoardID(ctx context.Context, projectID, boardID uuid.UUID) ([]*dto.BoardColumnView, error) {
	err := ensureProjectExists(ctx, uc.projectRepo, projectID)
	if err != nil {
		return nil, err
	}
	err = ensureBoardExists(ctx, uc.boardRepo, projectID, boardID)
	if err != nil {
		return nil, err
	}

	boardColumns, err := uc.boardColumnRepo.FindAllByBoardIDAndProjectID(ctx, projectID, boardID)
	if err != nil {
		return nil, fmt.Errorf("error fetching board columns: %w", err)
	}

	views := make([]*dto.BoardColumnView, len(boardColumns))
	for i := range boardColumns {
		views[i] = mapToBoardColumnView(&boardColumns[i], projectID)
	}

	return views, nil
}

func (uc *BoardColumnUseCase) Update(ctx context.Context, cmd dto.UpdateBoardColumnCommand) (*dto.BoardColumnView, error) {
	err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID)
	if err != nil {
		return nil, err
	}
	err = ensureBoardExists(ctx, uc.boardRepo, cmd.ProjectID, cmd.BoardID)
	if err != nil {
		return nil, err
	}

	bc, err := uc.boardColumnRepo.FindByIDAndBoardID(ctx, cmd.BoardID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching board column: %w", err)
	}
	if bc == nil {
		return nil, ErrBoardColumnNotFound
	}
	if cmd.Name != nil {
		bc.Name = *cmd.Name
	}
	if cmd.WipLimit != nil {
		bc.WipLimit = cmd.WipLimit
	}
	now := time.Now()
	bc.UpdatedAt = &now

	saved, err := uc.boardColumnRepo.Save(ctx, bc)
	if err != nil {
		return nil, err
	}
	return mapToBoardColumnView(saved, cmd.ProjectID), nil
}

func (uc *BoardColumnUseCase) Delete(ctx context.Context, projectID, boardID, id uuid.UUID) error {
	err := ensureProjectExists(ctx, uc.projectRepo, projectID)
	if err != nil {
		return err
	}
	err = ensureBoardExists(ctx, uc.boardRepo, projectID, boardID)
	if err != nil {
		return err
	}
	deleted, err := uc.boardColumnRepo.DeleteByIDAndBoardID(ctx, boardID, id)
	if err != nil {
		return fmt.Errorf("error deleting board column: %w", err)
	}
	if !deleted {
		return ErrBoardColumnNotFound
	}
	return nil
}

func (uc *BoardColumnUseCase) Reorder(ctx context.Context, cmd dto.ReorderBoardColumnsCommand) ([]*dto.BoardColumnView, error) {
	err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID)
	if err != nil {
		return nil, err
	}
	err = ensureBoardExists(ctx, uc.boardRepo, cmd.ProjectID, cmd.BoardID)
	if err != nil {
		return nil, err
	}

	existingCols, err := uc.boardColumnRepo.FindAllByBoardIDAndProjectID(ctx, cmd.ProjectID, cmd.BoardID)
	if err != nil {
		return nil, fmt.Errorf("error fetching board columns: %w", err)
	}
	if len(existingCols) != len(cmd.Positions) {
		return nil, fmt.Errorf("column count mismatch: expected %d, got %d", len(existingCols), len(cmd.Positions))
	}
	existingMap := make(map[uuid.UUID]bool, len(existingCols))
	for _, col := range existingCols {
		existingMap[col.ID] = true
	}
	for _, posID := range cmd.Positions {
		if !existingMap[posID] {
			return nil, ErrBoardColumnNotFound
		}
	}

	err = uc.boardColumnRepo.Reorder(ctx, cmd.BoardID, cmd.Positions)
	if err != nil {
		return nil, fmt.Errorf("error reordering board columns: %w", err)
	}

	return uc.GetAllByBoardID(ctx, cmd.ProjectID, cmd.BoardID)
}

func mapToBoardColumnView(bc *model.BoardColumn, projectID uuid.UUID) *dto.BoardColumnView {
	return &dto.BoardColumnView{
		ID:        bc.ID,
		ProjectID: projectID,
		BoardID:   bc.BoardID,
		Name:      bc.Name,
		Position:  bc.Position,
		WipLimit:  bc.WipLimit,
		CreatedAt: bc.CreatedAt,
		UpdatedAt: bc.UpdatedAt,
	}
}

func ensureBoardColumnExists(ctx context.Context, bcRepo port.BoardColumnRepository, projectID, boardID, bcID uuid.UUID) error {
	exists, err := bcRepo.ExistsByIDAndBoardIDAndProjectID(ctx, bcID, boardID, projectID)
	if err != nil {
		return fmt.Errorf("error checking if board column exists: %w", err)
	}
	if !exists {
		return ErrBoardColumnNotFound
	}
	return nil
}
