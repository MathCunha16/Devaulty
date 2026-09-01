package mcp

import (
	"context"
	"testing"

	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/dto"

	"github.com/stretchr/testify/require"
)

func TestProblemTools_CreateListGetUpdateDeleteAndStatus(t *testing.T) {
	app := setupMCPTestApp(t)
	project := app.createProject(t, "Problem project")
	tools := NewProblemTools(app.problemUseCase)

	createResult, err := tools.create(context.Background(), newToolRequest("create_problem", map[string]any{
		"projectID":        project.ID.String(),
		"title":            "Gateway timeout",
		"errorDescription": "API timed out while proxying",
		"solution":         "Add a retry policy",
		"status":           string(model.ProblemStatusOpen),
		"severity":         string(model.ProblemSeverityHigh),
	}))
	require.NoError(t, err)
	require.False(t, createResult.IsError)

	created, ok := createResult.StructuredContent.(*dto.ProblemView)
	require.True(t, ok)
	require.Equal(t, "Gateway timeout", created.Title)

	listResult, err := tools.list(context.Background(), newToolRequest("list_problems", map[string]any{
		"projectID":  project.ID.String(),
		"pageNumber": 1,
		"pageSize":   10,
	}))
	require.NoError(t, err)
	require.False(t, listResult.IsError)
	require.NotNil(t, listResult.StructuredContent)

	getResult, err := tools.get(context.Background(), newToolRequest("get_problem", map[string]any{
		"projectID": project.ID.String(),
		"id":        created.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, getResult.IsError)

	fetched, ok := getResult.StructuredContent.(*dto.ProblemView)
	require.True(t, ok)
	require.Equal(t, created.ID, fetched.ID)

	updateResult, err := tools.update(context.Background(), newToolRequest("update_problem", map[string]any{
		"projectID":        project.ID.String(),
		"id":               created.ID.String(),
		"title":            "Gateway timeout after deploy",
		"errorDescription": "API timed out after the rollout",
		"solution":         "Set a larger timeout and add retries",
		"severity":         string(model.ProblemSeverityCritical),
	}))
	require.NoError(t, err)
	require.False(t, updateResult.IsError)

	updated, ok := updateResult.StructuredContent.(*dto.ProblemView)
	require.True(t, ok)
	require.Equal(t, "Gateway timeout after deploy", updated.Title)

	statusResult, err := tools.updateStatus(context.Background(), newToolRequest("update_problem_status", map[string]any{
		"projectID": project.ID.String(),
		"id":        created.ID.String(),
		"status":    string(model.ProblemStatusWorkingOnIt),
	}))
	require.NoError(t, err)
	require.False(t, statusResult.IsError)

	statusUpdated, ok := statusResult.StructuredContent.(*dto.ProblemView)
	require.True(t, ok)
	require.Equal(t, model.ProblemStatusWorkingOnIt, statusUpdated.Status)

	deleteResult, err := tools.delete(context.Background(), newToolRequest("delete_problem", map[string]any{
		"projectID": project.ID.String(),
		"id":        created.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, deleteResult.IsError)
	require.Equal(t, "Problem Deleted!", textFromResult(t, deleteResult))

	invalidResult, err := tools.get(context.Background(), newToolRequest("get_problem", map[string]any{
		"projectID": project.ID.String(),
		"id":        "not-a-uuid",
	}))
	require.NoError(t, err)
	require.True(t, invalidResult.IsError)
	require.Contains(t, textFromResult(t, invalidResult), "invalid UUID format")
}
