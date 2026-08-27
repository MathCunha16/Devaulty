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
	ErrCardNotFound                     = errors.New("card not found")
	ErrNotPossibleToGetCardNextPosition = errors.New("not possible to get card next position")
)

type CardUseCase struct {
	cardRepo        port.CardRepository
	boardRepo       port.BoardRepository
	boardColumnRepo port.BoardColumnRepository
	projectRepo     port.ProjectRepository
	itemTagRepo     port.ItemTagRepository
}

func NewCardUseCase(cardRepo port.CardRepository, boardRepo port.BoardRepository, boardColumnRepo port.BoardColumnRepository,
	projectRepo port.ProjectRepository, itemTagRepo port.ItemTagRepository) *CardUseCase {
	return &CardUseCase{
		cardRepo:        cardRepo,
		boardRepo:       boardRepo,
		boardColumnRepo: boardColumnRepo,
		projectRepo:     projectRepo,
		itemTagRepo:     itemTagRepo,
	}
}

func (uc *CardUseCase) Create(ctx context.Context, cmd dto.CreateCardCommand) (*dto.CardView, error) {
	err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID)
	if err != nil {
		return nil, err
	}
	err = ensureBoardExists(ctx, uc.boardRepo, cmd.ProjectID, cmd.BoardID)
	if err != nil {
		return nil, err
	}
	err = ensureBoardColumnExists(ctx, uc.boardColumnRepo, cmd.ProjectID, cmd.BoardID, cmd.ColumnID)
	if err != nil {
		return nil, err
	}

	nextPosition, err := uc.cardRepo.GetNextPosition(ctx, cmd.ColumnID)
	if err != nil {
		return nil, ErrNotPossibleToGetCardNextPosition
	}

	card := model.Card{
		ID:          uuid.New(),
		ColumnID:    cmd.ColumnID,
		Title:       cmd.Title,
		Description: cmd.Description,
		Position:    nextPosition,
		Priority:    cmd.Priority,
		DueDate:     cmd.DueDate,
		LinkedItems: nil,
		BaseEntity: model.BaseEntity{
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		},
	}

	saved, err := uc.cardRepo.Save(ctx, &card)
	if err != nil {
		return nil, err
	}

	var linkedItems []model.CardItem
	if len(cmd.LinkedItems) > 0 {
		linkedItems = make([]model.CardItem, len(cmd.LinkedItems))
		for i, item := range cmd.LinkedItems {
			linkedItems[i] = model.CardItem{
				CardID:   saved.ID,
				ItemID:   item.ItemID,
				ItemType: item.ItemType,
				BaseEntity: model.BaseEntity{
					CreatedAt: time.Now(),
					UpdatedAt: nil,
				},
			}
		}

		if err := uc.cardRepo.SaveLinkedItems(ctx, saved.ID, linkedItems); err != nil {
			return nil, fmt.Errorf("error saving linked items: %w", err)
		}
		saved.LinkedItems = linkedItems
	}
	return mapToCardView(saved, nil, cmd.ProjectID, cmd.BoardID), nil
}

func (uc *CardUseCase) GetByID(ctx context.Context, projectID, boardID, id uuid.UUID) (*dto.CardView, error) {
	err := ensureProjectExists(ctx, uc.projectRepo, projectID)
	if err != nil {
		return nil, err
	}
	err = ensureBoardExists(ctx, uc.boardRepo, projectID, boardID)
	if err != nil {
		return nil, err
	}
	card, err := uc.cardRepo.FindByIDAndBoardID(ctx, boardID, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching card: %w", err)
	}
	if card == nil {
		return nil, ErrCardNotFound
	}
	linkedItems, err := uc.cardRepo.FindLinkedItemsByCardID(ctx, card.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching linked items: %w", err)
	}
	card.LinkedItems = linkedItems
	tags, err := uc.itemTagRepo.FindTagsForItem(ctx, model.ItemTypeCard, projectID, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for card: %w", err)
	}
	return mapToCardView(card, tags, projectID, boardID), nil
}

func (uc *CardUseCase) GetAllByBoardID(ctx context.Context, projectID, boardID uuid.UUID) ([]*dto.CardSummaryView, error) {
	err := ensureProjectExists(ctx, uc.projectRepo, projectID)
	if err != nil {
		return nil, err
	}
	err = ensureBoardExists(ctx, uc.boardRepo, projectID, boardID)
	if err != nil {
		return nil, err
	}

	cards, err := uc.cardRepo.FindAllByBoardIDAndProjectID(ctx, projectID, boardID)
	if err != nil {
		return nil, fmt.Errorf("error fetching cards: %w", err)
	}
	if len(cards) <= 0 {
		return []*dto.CardSummaryView{}, nil
	}

	cardIDs := make([]uuid.UUID, len(cards))
	for i, c := range cards {
		cardIDs[i] = c.ID
	}
	tagsMap, err := uc.itemTagRepo.FindTagsForItems(ctx, model.ItemTypeCard, projectID, cardIDs)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for cards: %w", err)
	}

	views := make([]*dto.CardSummaryView, len(cards))
	for i, c := range cards {
		tags := tagsMap[c.ID]
		tagSummaries := make([]dto.TagSummary, len(tags))
		for j, t := range tags {
			tagSummaries[j] = dto.TagSummary{
				ID:    t.ID,
				Name:  t.Name,
				Color: t.Color,
			}
		}
		views[i] = &dto.CardSummaryView{
			ID:        c.ID,
			ProjectID: projectID,
			BoardID:   boardID,
			ColumnID:  c.ColumnID,
			Title:     c.Title,
			Position:  c.Position,
			Priority:  c.Priority,
			DueDate:   c.DueDate,
			Tags:      tagSummaries,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
	}
	return views, nil
}

func (uc *CardUseCase) Update(ctx context.Context, cmd dto.UpdateCardCommand) (*dto.CardView, error) {
	err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID)
	if err != nil {
		return nil, err
	}
	err = ensureBoardExists(ctx, uc.boardRepo, cmd.ProjectID, cmd.BoardID)
	if err != nil {
		return nil, err
	}

	card, err := uc.cardRepo.FindByIDAndBoardID(ctx, cmd.BoardID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching card: %w", err)
	}
	if card == nil {
		return nil, ErrCardNotFound
	}
	if cmd.Title != nil {
		card.Title = *cmd.Title
	}
	if cmd.Description != nil {
		card.Description = cmd.Description
	}
	if cmd.Priority != nil {
		card.Priority = cmd.Priority
	}
	if cmd.DueDate != nil {
		card.DueDate = cmd.DueDate
	}
	now := time.Now()
	card.UpdatedAt = &now

	saved, err := uc.cardRepo.Save(ctx, card)
	if err != nil {
		return nil, err
	}

	if cmd.LinkedItems != nil && len(cmd.LinkedItems) > 0 {
		linkedItems := make([]model.CardItem, len(cmd.LinkedItems))
		for i, item := range cmd.LinkedItems {
			linkedItems[i] = model.CardItem{
				CardID:   saved.ID,
				ItemID:   item.ItemID,
				ItemType: item.ItemType,
				BaseEntity: model.BaseEntity{
					CreatedAt: time.Now(),
					UpdatedAt: nil,
				},
			}
		}

		if err := uc.cardRepo.SaveLinkedItems(ctx, saved.ID, linkedItems); err != nil {
			return nil, fmt.Errorf("error saving linked items: %w", err)
		}
		saved.LinkedItems = linkedItems
	} else {
		linkedItems, err := uc.cardRepo.FindLinkedItemsByCardID(ctx, saved.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching linked items: %w", err)
		}
		saved.LinkedItems = linkedItems
	}
	tags, err := uc.itemTagRepo.FindTagsForItem(ctx, model.ItemTypeCard, cmd.ProjectID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for card: %w", err)
	}
	return mapToCardView(saved, tags, cmd.ProjectID, cmd.BoardID), nil
}

func (uc *CardUseCase) Delete(ctx context.Context, projectID, boardID, id uuid.UUID) error {
	err := ensureProjectExists(ctx, uc.projectRepo, projectID)
	if err != nil {
		return err
	}
	err = ensureBoardExists(ctx, uc.boardRepo, projectID, boardID)
	if err != nil {
		return err
	}
	deleted, err := uc.cardRepo.DeleteByIDAndBoardID(ctx, boardID, id)
	if err != nil {
		return fmt.Errorf("error deleting card: %w", err)
	}
	if !deleted {
		return ErrCardNotFound
	}
	if err := uc.itemTagRepo.RemoveAllTagsFromItem(ctx, model.ItemTypeCard, id); err != nil {
		log.Printf("warning: failed to remove tags from card %s: %v", id, err)
	}
	return nil
}

func (uc *CardUseCase) Move(ctx context.Context, cmd dto.MoveCardCommand) error {
	err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID)
	if err != nil {
		return err
	}
	err = ensureBoardExists(ctx, uc.boardRepo, cmd.ProjectID, cmd.BoardID)
	if err != nil {
		return err
	}
	err = ensureBoardColumnExists(ctx, uc.boardColumnRepo, cmd.ProjectID, cmd.BoardID, cmd.TargetColumnID)
	if err != nil {
		return err
	}

	card, err := uc.cardRepo.FindByIDAndBoardID(ctx, cmd.BoardID, cmd.ID)
	if err != nil {
		return fmt.Errorf("error fetching card: %w", err)
	}
	if card == nil {
		return ErrCardNotFound
	}
	err = uc.cardRepo.MoveCard(ctx, cmd.ID, card.ColumnID, cmd.TargetColumnID, *cmd.Position)
	if err != nil {
		return fmt.Errorf("error moving card: %w", err)
	}
	return nil
}

func mapToCardView(card *model.Card, tags []model.Tag, projectID, boardID uuid.UUID) *dto.CardView {
	tagSummaries := make([]dto.TagSummary, len(tags))
	for i, t := range tags {
		tagSummaries[i] = dto.TagSummary{
			ID:    t.ID,
			Name:  t.Name,
			Color: t.Color,
		}
	}
	return &dto.CardView{
		ID:          card.ID,
		ProjectID:   projectID,
		BoardID:     boardID,
		ColumnID:    card.ColumnID,
		Title:       card.Title,
		Description: card.Description,
		Position:    card.Position,
		Priority:    card.Priority,
		DueDate:     card.DueDate,
		LinkedItems: card.LinkedItems,
		Tags:        tagSummaries,
		CreatedAt:   card.CreatedAt,
		UpdatedAt:   card.UpdatedAt,
	}
}
