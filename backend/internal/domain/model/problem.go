package model

import "github.com/google/uuid"

type ProblemStatus string

const (
	ProblemStatusOpen        ProblemStatus = "OPEN"
	ProblemStatusWorkingOnIt ProblemStatus = "WORKING_ON"
	ProblemStatusResolved    ProblemStatus = "RESOLVED"
	ProblemStatusWontFix     ProblemStatus = "WONT_FIX"
)

var ProblemStatuses = []ProblemStatus{ProblemStatusOpen, ProblemStatusWorkingOnIt, ProblemStatusResolved, ProblemStatusWontFix}

type ProblemSeverity string

const (
	ProblemSeverityLow      ProblemSeverity = "LOW"
	ProblemSeverityMedium   ProblemSeverity = "MEDIUM"
	ProblemSeverityHigh     ProblemSeverity = "HIGH"
	ProblemSeverityCritical ProblemSeverity = "CRITICAL"
)

var ProblemSeverities = []ProblemSeverity{ProblemSeverityLow, ProblemSeverityMedium, ProblemSeverityHigh, ProblemSeverityCritical}

type Problem struct {
	ID               uuid.UUID       `json:"id" db:"id"`
	ProjectID        uuid.UUID       `json:"projectId" db:"project_id"`
	Title            string          `json:"title" db:"title"`
	ErrorDescription string          `json:"errorDescription" db:"error_description"`
	Solution         *string         `json:"solution,omitempty" db:"solution"`
	Status           ProblemStatus   `json:"status" db:"status"`
	Severity         ProblemSeverity `json:"severity" db:"severity"`
	BaseEntity
}
