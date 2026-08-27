package usecase

import (
	"context"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/domain/port"
	"devaulty-backend/internal/dto"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

var (
	ErrBoardNotFound                  = errors.New("board not found")
	ErrNotPossibleToUnsetDefaultBoard = errors.New("not possible to unset default board")
)

type BoardUseCase struct {
	boardRepo       port.BoardRepository
	boardColumnRepo port.BoardColumnRepository
	projectRepo     port.ProjectRepository
	itemTagRepo     port.ItemTagRepository
}

func NewBoardUseCase(boardRepo port.BoardRepository, boardColumnRepo port.BoardColumnRepository, projectRepo port.ProjectRepository, itemTagRepo port.ItemTagRepository) *BoardUseCase {
	return &BoardUseCase{
		boardRepo:       boardRepo,
		boardColumnRepo: boardColumnRepo,
		projectRepo:     projectRepo,
		itemTagRepo:     itemTagRepo,
	}
}

func (uc *BoardUseCase) Create(ctx context.Context, cmd dto.CreateBoardCommand) (*dto.BoardView, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}

	if cmd.IsDefault {
		result, err := uc.boardRepo.UnsetAllDefaultsByProjectID(ctx, cmd.ProjectID)
		if err != nil {
			return nil, err
		}
		if !result {
			return nil, ErrNotPossibleToUnsetDefaultBoard
		}

	} else {
		defaultBoard, err := uc.boardRepo.FindDefaultByProjectID(ctx, cmd.ProjectID)
		if err != nil {
			return nil, err
		}
		if defaultBoard == nil {
			cmd.IsDefault = true
		}
	}

	board := model.Board{
		ID:          uuid.New(),
		ProjectID:   cmd.ProjectID,
		Name:        cmd.Name,
		Description: cmd.Description,
		IsDefault:   cmd.IsDefault,
		BaseEntity: model.BaseEntity{
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		},
	}

	saved, err := uc.boardRepo.Save(ctx, &board)
	if err != nil {
		return nil, err
	}

	err = uc.createDefaultBoardColumns(ctx, saved.ID)
	if err != nil {
		_, _ = uc.boardRepo.DeleteByIDAndProjectID(ctx, cmd.ProjectID, saved.ID)
		return nil, fmt.Errorf("error creating default board columns: %w", err)
	}

	return mapBoardToView(saved, nil), nil
}

func (uc *BoardUseCase) GetByID(ctx context.Context, projectID, id uuid.UUID) (*dto.BoardView, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return nil, err
	}
	board, err := uc.boardRepo.FindByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching board: %w", err)
	}
	if board == nil {
		return nil, ErrBoardNotFound
	}
	tags, err := uc.itemTagRepo.FindTagsForItem(ctx, model.ItemTypeBoard, projectID, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for board: %w", err)
	}
	return mapBoardToView(board, tags), nil
}

func (uc *BoardUseCase) GetDefaultByProjectID(ctx context.Context, projectID uuid.UUID) (*dto.BoardView, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return nil, err
	}
	board, err := uc.boardRepo.FindDefaultByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if board == nil {
		return nil, ErrBoardNotFound
	}
	tags, err := uc.itemTagRepo.FindTagsForItem(ctx, model.ItemTypeBoard, projectID, board.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for board: %w", err)
	}
	return mapBoardToView(board, tags), nil
}

func (uc *BoardUseCase) GetAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[dto.BoardView], error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return model.Page[dto.BoardView]{}, err
	}
	boardPage, err := uc.boardRepo.FindAllByProjectID(ctx, projectID, page, size)
	if err != nil {
		return model.Page[dto.BoardView]{}, err
	}
	if len(boardPage.Content) == 0 {
		return model.NewPage([]dto.BoardView{}, boardPage.Number, boardPage.Size, boardPage.TotalElements), nil
	}

	boardIDs := make([]uuid.UUID, len(boardPage.Content))
	for i, b := range boardPage.Content {
		boardIDs[i] = b.ID
	}

	tagsMap, err := uc.itemTagRepo.FindTagsForItems(ctx, model.ItemTypeBoard, projectID, boardIDs)
	if err != nil {
		return model.Page[dto.BoardView]{}, fmt.Errorf("error fetching tags for boards: %w", err)
	}

	views := make([]dto.BoardView, len(boardPage.Content))
	for i, b := range boardPage.Content {
		tags := tagsMap[b.ID]
		views[i] = *mapBoardToView(&b, tags)
	}

	return model.NewPage(views, boardPage.Number, boardPage.Size, boardPage.TotalElements), nil
}

func (uc *BoardUseCase) Update(ctx context.Context, cmd dto.UpdateBoardCommand) (*dto.BoardView, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}

	board, err := uc.boardRepo.FindByIDAndProjectID(ctx, cmd.ProjectID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching board: %w", err)
	}
	if board == nil {
		return nil, ErrBoardNotFound
	}

	if cmd.Name != nil {
		board.Name = *cmd.Name
	}
	if cmd.Description != nil {
		board.Description = cmd.Description
	}
	if cmd.IsDefault != nil {
		if *cmd.IsDefault {
			result, err := uc.boardRepo.UnsetAllDefaultsByProjectID(ctx, cmd.ProjectID)
			if err != nil {
				return nil, err
			}
			if !result {
				return nil, ErrNotPossibleToUnsetDefaultBoard
			}
			board.IsDefault = *cmd.IsDefault
		} else {
			board.IsDefault = *cmd.IsDefault
		}
	}
	now := time.Now()
	board.UpdatedAt = &now
	saved, err := uc.boardRepo.Save(ctx, board)
	if err != nil {
		return nil, err
	}

	tags, err := uc.itemTagRepo.FindTagsForItem(ctx, model.ItemTypeBoard, cmd.ProjectID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for board: %w", err)
	}

	return mapBoardToView(saved, tags), nil
}

func (uc *BoardUseCase) Delete(ctx context.Context, projectID, id uuid.UUID) error {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return err
	}
	deleted, err := uc.boardRepo.DeleteByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return fmt.Errorf("error deleting board: %w", err)
	}
	if !deleted {
		return ErrBoardNotFound
	}

	if err := uc.itemTagRepo.RemoveAllTagsFromItem(ctx, model.ItemTypeBoard, id); err != nil {
		log.Printf("warning: failed to remove tags from board %s: %v", id, err)
	}
	return nil
}

func mapBoardToView(board *model.Board, tags []model.Tag) *dto.BoardView {
	tagSummaries := make([]dto.TagSummary, len(tags))
	for i, t := range tags {
		tagSummaries[i] = dto.TagSummary{
			ID:    t.ID,
			Name:  t.Name,
			Color: t.Color,
		}
	}
	return &dto.BoardView{
		ID:          board.ID,
		ProjectID:   board.ProjectID,
		Name:        board.Name,
		Description: board.Description,
		IsDefault:   board.IsDefault,
		Tags:        tagSummaries,
		CreatedAt:   board.CreatedAt,
		UpdatedAt:   board.UpdatedAt,
	}
}

func (uc *BoardUseCase) createDefaultBoardColumns(ctx context.Context, boardID uuid.UUID) error {
	defaultColumns := []string{"Backlog", "In Progress", "Review", "Done"}
	for pos, name := range defaultColumns {
		col := model.BoardColumn{
			ID:       uuid.New(),
			BoardID:  boardID,
			Name:     name,
			Position: uint8(pos),
			BaseEntity: model.BaseEntity{
				CreatedAt: time.Now(),
				UpdatedAt: nil,
			},
		}
		switch name {
		case "In Progress", "Review":
			wipLimit := uint16(3)
			col.WipLimit = &wipLimit
		default:
			col.WipLimit = nil
		}
		if _, err := uc.boardColumnRepo.Save(ctx, &col); err != nil {
			return err
		}
	}
	return nil
}

func ensureBoardExists(ctx context.Context, boardRepo port.BoardRepository, projectID, boardID uuid.UUID) error {
	exits, err := boardRepo.ExistsByIDAndProjectID(ctx, boardID, projectID)
	if err != nil {
		return fmt.Errorf("error checking if board exists: %w", err)
	}
	if !exits {
		return ErrBoardNotFound
	}
	return nil
}
