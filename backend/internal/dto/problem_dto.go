package dto

import (
	"time"

	"devaulty-backend/internal/domain/model"

	"github.com/google/uuid"
)

type CreateProblemCommand struct {
	ProjectID        uuid.UUID             `json:"projectID"`
	Title            string                `json:"title" binding:"required,min=2,max=255"`
	ErrorDescription string                `json:"errorDescription" binding:"required,min=2,max=255"`
	Solution         *string               `json:"solution,omitempty" binding:"omitempty,min=2,max=255"`
	Status           model.ProblemStatus   `json:"status" binding:"required"`
	Severity         model.ProblemSeverity `json:"severity" binding:"required"`
}

type UpdateProblemCommand struct {
	ProjectID        uuid.UUID              `json:"projectID"`
	ID               uuid.UUID              `json:"id"`
	Title            *string                `json:"title,omitempty" binding:"omitempty,min=2,max=255"`
	ErrorDescription *string                `json:"errorDescription,omitempty" binding:"omitempty,min=2,max=255"`
	Solution         *string                `json:"solution,omitempty" binding:"omitempty,min=2,max=255"`
	Severity         *model.ProblemSeverity `json:"severity,omitempty" binding:"omitempty"`
}

type UpdateProblemStatusCommand struct {
	ProjectID uuid.UUID           `json:"projectID"`
	ID        uuid.UUID           `json:"id"`
	Status    model.ProblemStatus `json:"status" binding:"required"`
}

type ProblemSummary struct {
	ID        uuid.UUID             `json:"id"`
	ProjectID uuid.UUID             `json:"projectId"`
	Title     string                `json:"title"`
	Status    model.ProblemStatus   `json:"status"`
	Severity  model.ProblemSeverity `json:"severity"`
	Tags      []TagSummary          `json:"tags"`
	CreatedAt *time.Time            `json:"createdAt,omitempty"`
	UpdatedAt *time.Time            `json:"updatedAt,omitempty"`
}

type ProblemView struct {
	ID               uuid.UUID             `json:"id"`
	ProjectID        uuid.UUID             `json:"projectId"`
	Title            string                `json:"title"`
	ErrorDescription string                `json:"errorDescription"`
	Solution         *string               `json:"solution,omitempty"`
	Status           model.ProblemStatus   `json:"status"`
	Severity         model.ProblemSeverity `json:"severity"`
	Tags             []TagSummary          `json:"tags"`
	CreatedAt        *time.Time            `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time            `json:"updatedAt,omitempty"`
}
