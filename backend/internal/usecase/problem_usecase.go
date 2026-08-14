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
	ErrProblemNotFound = errors.New("problem not found")
)

type ProblemUseCase struct {
	problemRepo port.ProblemRepository
	projectRepo port.ProjectRepository
	itemTagRepo port.ItemTagRepository
}

func NewProblemUseCase(problemRepo port.ProblemRepository, projectRepo port.ProjectRepository, itemTagRepo port.ItemTagRepository) *ProblemUseCase {
	return &ProblemUseCase{
		problemRepo: problemRepo,
		projectRepo: projectRepo,
		itemTagRepo: itemTagRepo,
	}
}

func (uc *ProblemUseCase) Create(ctx context.Context, cmd dto.CreateProblemCommand) (*dto.ProblemView, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}

	problem := model.Problem{
		ID:               uuid.New(),
		ProjectID:        cmd.ProjectID,
		Title:            cmd.Title,
		ErrorDescription: cmd.ErrorDescription,
		Solution:         cmd.Solution,
		Status:           cmd.Status,
		Severity:         cmd.Severity,
		BaseEntity: model.BaseEntity{
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		},
	}

	saved, err := uc.problemRepo.Save(ctx, &problem)
	if err != nil {
		return nil, err
	}
	return mapProblemToView(saved, nil), nil
}

func (uc *ProblemUseCase) GetByID(ctx context.Context, projectID, id uuid.UUID) (*dto.ProblemView, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return nil, err
	}
	problem, err := uc.problemRepo.FindByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return nil, fmt.Errorf("error trying to find problem: %w", err)
	}
	if problem == nil {
		return nil, ErrProblemNotFound
	}
	tags, err := uc.itemTagRepo.FindTagsForItem(ctx, model.ItemTypeProblem, projectID, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for problem: %w", err)
	}
	return mapProblemToView(problem, tags), nil
}

func (uc *ProblemUseCase) GetAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[dto.ProblemSummary], error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return model.Page[dto.ProblemSummary]{}, err
	}
	problemPage, err := uc.problemRepo.FindAllByProjectID(ctx, projectID, page, size)
	if err != nil {
		return model.Page[dto.ProblemSummary]{}, err
	}
	if len(problemPage.Content) == 0 {
		return model.NewPage([]dto.ProblemSummary{}, problemPage.Number, problemPage.Size, problemPage.TotalElements), nil
	}

	problemIDs := make([]uuid.UUID, len(problemPage.Content))
	for i, p := range problemPage.Content {
		problemIDs[i] = p.ID
	}

	tagsMap, err := uc.itemTagRepo.FindTagsForItems(ctx, model.ItemTypeProblem, projectID, problemIDs)
	if err != nil {
		return model.Page[dto.ProblemSummary]{}, fmt.Errorf("error fetching tags for problems: %w", err)
	}

	summaries := make([]dto.ProblemSummary, len(problemPage.Content))
	for i, p := range problemPage.Content {
		tags := tagsMap[p.ID]
		tagSummaries := make([]dto.TagSummary, len(tags))
		for j, t := range tags {
			tagSummaries[j] = dto.TagSummary{
				ID:    t.ID,
				Name:  t.Name,
				Color: t.Color,
			}
		}
		summaries[i] = dto.ProblemSummary{
			ID:        p.ID,
			ProjectID: p.ProjectID,
			Title:     p.Title,
			Status:    p.Status,
			Severity:  p.Severity,
			Tags:      tagSummaries,
			CreatedAt: &p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}

	return model.NewPage(summaries, problemPage.Number, problemPage.Size, problemPage.TotalElements), nil
}

func (uc *ProblemUseCase) Update(ctx context.Context, cmd dto.UpdateProblemCommand) (*dto.ProblemView, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}
	problem, err := uc.problemRepo.FindByIDAndProjectID(ctx, cmd.ProjectID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error trying to find problem: %w", err)
	}
	if problem == nil {
		return nil, ErrProblemNotFound
	}

	if cmd.Title != nil {
		problem.Title = *cmd.Title
	}
	if cmd.ErrorDescription != nil {
		problem.ErrorDescription = *cmd.ErrorDescription
	}
	if cmd.Solution != nil {
		problem.Solution = cmd.Solution
	}
	if cmd.Severity != nil {
		problem.Severity = *cmd.Severity
	}
	now := time.Now()
	problem.UpdatedAt = &now
	saved, err := uc.problemRepo.Save(ctx, problem)
	if err != nil {
		return nil, err
	}
	tags, err := uc.itemTagRepo.FindTagsForItem(ctx, model.ItemTypeProblem, cmd.ProjectID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for problem: %w", err)
	}
	return mapProblemToView(saved, tags), nil
}

func (uc *ProblemUseCase) UpdateStatus(ctx context.Context, cmd dto.UpdateProblemStatusCommand) (*dto.ProblemView, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}
	problem, err := uc.problemRepo.FindByIDAndProjectID(ctx, cmd.ProjectID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error trying to find problem: %w", err)
	}
	if problem == nil {
		return nil, ErrProblemNotFound
	}

	problem.Status = cmd.Status
	now := time.Now()
	problem.UpdatedAt = &now
	saved, err := uc.problemRepo.Save(ctx, problem)
	if err != nil {
		return nil, err
	}
	tags, err := uc.itemTagRepo.FindTagsForItem(ctx, model.ItemTypeProblem, cmd.ProjectID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for problem: %w", err)
	}
	return mapProblemToView(saved, tags), nil
}

func (uc *ProblemUseCase) Delete(ctx context.Context, projectID, id uuid.UUID) error {
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return err
	}

	deleted, err := uc.problemRepo.DeleteByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return fmt.Errorf("error deleting problem: %w", err)
	}
	if !deleted {
		return ErrProblemNotFound
	}

	if err := uc.itemTagRepo.RemoveAllTagsFromItem(ctx, model.ItemTypeProblem, id); err != nil {
		log.Printf("warning: failed to remove tags from problem %s: %v", id, err)
	}
	return nil
}

func mapProblemToView(problem *model.Problem, tags []model.Tag) *dto.ProblemView {
	tagSummaries := make([]dto.TagSummary, len(tags))
	for i, t := range tags {
		tagSummaries[i] = dto.TagSummary{
			ID:    t.ID,
			Name:  t.Name,
			Color: t.Color,
		}
	}
	return &dto.ProblemView{
		ID:               problem.ID,
		ProjectID:        problem.ProjectID,
		Title:            problem.Title,
		ErrorDescription: problem.ErrorDescription,
		Solution:         problem.Solution,
		Status:           problem.Status,
		Severity:         problem.Severity,
		Tags:             tagSummaries,
		CreatedAt:        &problem.CreatedAt,
		UpdatedAt:        problem.UpdatedAt,
	}
}
