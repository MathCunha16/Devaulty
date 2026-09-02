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
	ErrNoteNotFound          = errors.New("note not found")
	ErrNoteAlreadyArchived   = errors.New("note already archived")
	ErrNoteAlreadyUnarchived = errors.New("note already unarchived")
)

type NoteUseCase struct {
	noteRepo    port.NoteRepository
	projectRepo port.ProjectRepository
	itemTagRepo port.ItemTagRepository
}

func NewNoteUseCase(noteRepo port.NoteRepository, projectRepo port.ProjectRepository, itemTagRepo port.ItemTagRepository) *NoteUseCase {
	return &NoteUseCase{
		noteRepo:    noteRepo,
		projectRepo: projectRepo,
		itemTagRepo: itemTagRepo,
	}
}

func (uc *NoteUseCase) Create(ctx context.Context, cmd dto.CreateNoteCommand) (*dto.NoteView, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}

	note := model.Note{
		ID:        uuid.New(),
		ProjectID: cmd.ProjectID,
		Title:     cmd.Title,
		Content:   &cmd.Content,
		Archived:  false,
		BaseEntity: model.BaseEntity{
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		},
	}
	saved, err := uc.noteRepo.Save(ctx, &note)
	if err != nil {
		return nil, err
	}
	return mapNoteToView(saved, nil), nil
}

func (uc *NoteUseCase) GetByID(ctx context.Context, projectID, id uuid.UUID) (*dto.NoteView, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return nil, err
	}
	note, err := uc.noteRepo.FindByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return nil, fmt.Errorf("error trying to find note: %w", err)
	}
	if note == nil {
		return nil, ErrNoteNotFound
	}
	tags, err := uc.itemTagRepo.FindTagsForItem(ctx, model.ItemTypeNote, projectID, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for note: %w", err)
	}
	return mapNoteToView(note, tags), nil
}

func (uc *NoteUseCase) GetAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[dto.NoteSummary], error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return model.Page[dto.NoteSummary]{}, err
	}

	notePage, err := uc.noteRepo.FindAllByProjectID(ctx, projectID, page, size)
	if err != nil {
		return model.Page[dto.NoteSummary]{}, err
	}
	if len(notePage.Content) == 0 {
		return model.NewPage([]dto.NoteSummary{}, notePage.Number, notePage.Size, notePage.TotalElements), nil
	}

	noteIDs := make([]uuid.UUID, len(notePage.Content))
	for i, p := range notePage.Content {
		noteIDs[i] = p.ID
	}

	notesMap, err := uc.itemTagRepo.FindTagsForItems(ctx, model.ItemTypeNote, projectID, noteIDs)
	if err != nil {
		return model.Page[dto.NoteSummary]{}, fmt.Errorf("error fetching tags for notes: %w", err)
	}

	summaries := make([]dto.NoteSummary, len(notePage.Content))
	for i, p := range notePage.Content {
		tags := notesMap[p.ID]
		tagSummaries := make([]dto.TagSummary, len(tags))
		for j, t := range tags {
			tagSummaries[j] = dto.TagSummary{
				ID:    t.ID,
				Name:  t.Name,
				Color: t.Color,
			}
		}
		summaries[i] = dto.NoteSummary{
			ID:        p.ID,
			ProjectID: p.ProjectID,
			Title:     p.Title,
			Archived:  p.Archived,
			Tags:      tagSummaries,
			CreatedAt: &p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}
	return model.NewPage(summaries, notePage.Number, notePage.Size, notePage.TotalElements), nil
}

func (uc *NoteUseCase) Update(ctx context.Context, cmd dto.UpdateNoteCommand) (*dto.NoteView, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}
	note, err := uc.noteRepo.FindByIDAndProjectID(ctx, cmd.ProjectID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error trying to find note: %w", err)
	}
	if note == nil {
		return nil, ErrNoteNotFound
	}

	if cmd.Title != nil {
		note.Title = *cmd.Title
	}
	if cmd.Content != nil {
		note.Content = cmd.Content
	}
	now := time.Now()
	note.UpdatedAt = &now
	savedNote, err := uc.noteRepo.Save(ctx, note)
	if err != nil {
		return nil, err
	}
	tags, err := uc.itemTagRepo.FindTagsForItem(ctx, model.ItemTypeNote, cmd.ProjectID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for note: %w", err)
	}

	return mapNoteToView(savedNote, tags), nil

}

func (uc *NoteUseCase) Delete(ctx context.Context, projectID, id uuid.UUID) error {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return err
	}
	deleted, err := uc.noteRepo.DeleteByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return fmt.Errorf("error deleting note: %w", err)
	}
	if !deleted {
		return ErrNoteNotFound
	}
	if err := uc.itemTagRepo.RemoveAllTagsFromItem(ctx, model.ItemTypeNote, id); err != nil {
		log.Printf("warning: failed to remove tags from note %s: %v", id, err)
	}
	return nil
}

func (uc *NoteUseCase) Archive(ctx context.Context, projectID, id uuid.UUID) error {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return err
	}
	note, err := uc.noteRepo.FindByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return fmt.Errorf("error trying to find note: %w", err)
	}
	if note == nil {
		return ErrNoteNotFound
	}
	if note.Archived {
		return ErrNoteAlreadyArchived
	}

	note.Archived = true
	now := time.Now()
	note.UpdatedAt = &now

	_, err = uc.noteRepo.Save(ctx, note)
	if err != nil {
		return fmt.Errorf("error archiving note: %w", err)
	}
	return nil
}

func (uc *NoteUseCase) Unarchive(ctx context.Context, projectID, id uuid.UUID) error {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return err
	}
	note, err := uc.noteRepo.FindByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return fmt.Errorf("error trying to find note: %w", err)
	}
	if note == nil {
		return ErrNoteNotFound
	}
	if !note.Archived {
		return ErrNoteAlreadyUnarchived
	}

	note.Archived = false
	now := time.Now()
	note.UpdatedAt = &now

	_, err = uc.noteRepo.Save(ctx, note)
	if err != nil {
		return fmt.Errorf("error unarchiving note: %w", err)
	}
	return nil
}

// --- auxiliary methods ---
func mapNoteToView(note *model.Note, tags []model.Tag) *dto.NoteView {
	tagSummaries := make([]dto.TagSummary, len(tags))
	for i, t := range tags {
		tagSummaries[i] = dto.TagSummary{
			ID:    t.ID,
			Name:  t.Name,
			Color: t.Color,
		}
	}
	var content string
	if note.Content != nil {
		content = *note.Content
	}
	return &dto.NoteView{
		ID:        note.ID,
		ProjectID: note.ProjectID,
		Title:     note.Title,
		Content:   content,
		Archived:  note.Archived,
		Tags:      tagSummaries,
		CreatedAt: &note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}
