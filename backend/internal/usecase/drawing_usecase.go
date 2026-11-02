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
	ErrDrawingNotFound = errors.New("drawing not found")
)

type DrawingUseCase struct {
	drawingRepo port.DrawingRepository
	projectRepo port.ProjectRepository
	itemTagRepo port.ItemTagRepository
}

func NewDrawingUseCase(drawingRepo port.DrawingRepository, projectRepo port.ProjectRepository, itemTagRepo port.ItemTagRepository) *DrawingUseCase {
	return &DrawingUseCase{
		drawingRepo: drawingRepo,
		projectRepo: projectRepo,
		itemTagRepo: itemTagRepo,
	}
}

func (uc *DrawingUseCase) Create(ctx context.Context, cmd dto.CreateDrawingCommand) (*dto.DrawingView, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}

	drawing := model.Drawing{
		ID:        uuid.New(),
		ProjectID: cmd.ProjectID,
		Title:     cmd.Title,
		Content:   nil,
		Thumbnail: nil,
		BaseEntity: model.BaseEntity{
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		},
	}

	saved, err := uc.drawingRepo.Save(ctx, &drawing)
	if err != nil {
		return nil, err
	}
	return mapDrawingToView(saved, nil), nil
}

func (uc *DrawingUseCase) GetByID(ctx context.Context, projectID, id uuid.UUID) (*dto.DrawingView, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return nil, err
	}

	drawing, err := uc.drawingRepo.FindByIDAndProjectID(ctx, id, projectID)
	if err != nil {
		return nil, fmt.Errorf("error trying to find drawing: %w", err)
	}
	if drawing == nil {
		return nil, ErrDrawingNotFound
	}

	tags, err := uc.itemTagRepo.FindTagsForItem(ctx, model.ItemTypeDrawing, projectID, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for drawing: %w", err)
	}

	return mapDrawingToView(drawing, tags), nil
}

func (uc *DrawingUseCase) GetAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[dto.DrawingSummary], error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return model.Page[dto.DrawingSummary]{}, err
	}

	drawingPage, err := uc.drawingRepo.FindAllByProjectID(ctx, projectID, page, size)
	if err != nil {
		return model.Page[dto.DrawingSummary]{}, err
	}
	if len(drawingPage.Content) == 0 {
		return model.NewPage([]dto.DrawingSummary{}, drawingPage.Number, drawingPage.Size, drawingPage.TotalElements), nil
	}

	drawingIDs := make([]uuid.UUID, len(drawingPage.Content))
	for i, p := range drawingPage.Content {
		drawingIDs[i] = p.ID
	}

	drawingsMap, err := uc.itemTagRepo.FindTagsForItems(ctx, model.ItemTypeDrawing, projectID, drawingIDs)
	if err != nil {
		return model.Page[dto.DrawingSummary]{}, fmt.Errorf("error fetching tags for drawings: %w", err)
	}

	summaries := make([]dto.DrawingSummary, len(drawingPage.Content))
	for i, p := range drawingPage.Content {
		tags := drawingsMap[p.ID]
		tagSummaries := make([]dto.TagSummary, len(tags))
		for j, t := range tags {
			tagSummaries[j] = dto.TagSummary{
				ID:    t.ID,
				Name:  t.Name,
				Color: t.Color,
			}
		}
		summaries[i] = dto.DrawingSummary{
			ID:        p.ID,
			ProjectID: p.ProjectID,
			Title:     p.Title,
			Thumbnail: p.Thumbnail,
			Tags:      tagSummaries,
			CreatedAt: &p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}

	return model.NewPage(summaries, drawingPage.Number, drawingPage.Size, drawingPage.TotalElements), nil
}

func (uc *DrawingUseCase) Update(ctx context.Context, cmd dto.UpdateDrawingCommand) (*dto.DrawingView, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}

	drawing, err := uc.drawingRepo.FindByIDAndProjectID(ctx, cmd.ID, cmd.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("error trying to find drawing: %w", err)
	}
	if drawing == nil {
		return nil, ErrDrawingNotFound
	}

	if cmd.Title != nil {
		drawing.Title = *cmd.Title
	}
	if cmd.Content != nil {
		drawing.Content = cmd.Content
	}
	if cmd.Thumbnail != nil {
		drawing.Thumbnail = cmd.Thumbnail
	}
	now := time.Now()
	drawing.UpdatedAt = &now
	savedDrawing, err := uc.drawingRepo.Save(ctx, drawing)
	if err != nil {
		return nil, fmt.Errorf("error trying to save drawing: %w", err)
	}
	tags, err := uc.itemTagRepo.FindTagsForItem(ctx, model.ItemTypeDrawing, cmd.ProjectID, savedDrawing.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for drawing: %w", err)
	}

	return mapDrawingToView(savedDrawing, tags), nil
}

func (uc *DrawingUseCase) Delete(ctx context.Context, projectID, id uuid.UUID) error {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return err
	}
	deleted, err := uc.drawingRepo.DeleteByIDAndProjectID(ctx, id, projectID)
	if err != nil {
		return fmt.Errorf("error deleting drawing: %w", err)
	}
	if !deleted {
		return ErrDrawingNotFound
	}
	if err := uc.itemTagRepo.RemoveAllTagsFromItem(ctx, model.ItemTypeDrawing, id); err != nil {
		return fmt.Errorf("error removing tags from drawing: %w", err)
	}

	return nil
}

// --- auxiliary methods ---
func mapDrawingToView(drawing *model.Drawing, tags []model.Tag) *dto.DrawingView {
	tagSummaries := make([]dto.TagSummary, len(tags))
	for i, t := range tags {
		tagSummaries[i] = dto.TagSummary{
			ID:    t.ID,
			Name:  t.Name,
			Color: t.Color,
		}
	}

	return &dto.DrawingView{
		ID:        drawing.ID,
		ProjectID: drawing.ProjectID,
		Title:     drawing.Title,
		Content:   drawing.Content,
		Thumbnail: drawing.Thumbnail,
		Tags:      tagSummaries,
		CreatedAt: &drawing.CreatedAt,
		UpdatedAt: drawing.UpdatedAt,
	}
}
