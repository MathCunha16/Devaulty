package mcp

import (
	"context"
	"testing"

	"devaulty-backend/internal/domain/model"

	"github.com/stretchr/testify/require"
)

func TestProjectTools_CreateAndGet(t *testing.T) {
	app := setupMCPTestApp(t)
	tools := NewProjectTools(app.projectUseCase)

	result, err := tools.create(context.Background(), newToolRequest("create_project", map[string]any{
		"name":        "Demo project",
		"description": "A project used for MCP tests",
		"color":       "#10b981",
		"icon":        "Folder",
	}))
	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.IsError)

	created, ok := result.StructuredContent.(*model.Project)
	require.True(t, ok)
	require.Equal(t, "Demo project", created.Name)

	getResult, err := tools.get(context.Background(), newToolRequest("get_project", map[string]any{
		"project_id": created.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, getResult.IsError)

	fetched, ok := getResult.StructuredContent.(*model.Project)
	require.True(t, ok)
	require.Equal(t, created.ID, fetched.ID)
	require.Equal(t, created.Name, fetched.Name)

	invalidResult, err := tools.get(context.Background(), newToolRequest("get_project", map[string]any{
		"project_id": "invalid-uuid",
	}))
	require.NoError(t, err)
	require.True(t, invalidResult.IsError)
	require.Contains(t, textFromResult(t, invalidResult), "invalid UUID format")
}

func TestProjectTools_UpdateArchiveUnarchiveAndDelete(t *testing.T) {
	app := setupMCPTestApp(t)
	tools := NewProjectTools(app.projectUseCase)
	project := app.createProject(t, "Project to update")

	updateResult, err := tools.update(context.Background(), newToolRequest("update_project", map[string]any{
		"id":          project.ID.String(),
		"name":        "Updated project name",
		"description": "Updated project description",
		"color":       "#3b82f6",
		"icon":        "Terminal",
	}))
	require.NoError(t, err)
	require.False(t, updateResult.IsError)

	updated, ok := updateResult.StructuredContent.(*model.Project)
	require.True(t, ok)
	require.Equal(t, "Updated project name", updated.Name)

	archiveResult, err := tools.archive(context.Background(), newToolRequest("archive_project", map[string]any{
		"project_id": project.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, archiveResult.IsError)
	require.Equal(t, "Project Archived!", textFromResult(t, archiveResult))

	unarchiveResult, err := tools.unarchive(context.Background(), newToolRequest("unarchive_project", map[string]any{
		"project_id": project.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, unarchiveResult.IsError)
	require.Equal(t, "Project Unarchived!", textFromResult(t, unarchiveResult))

	deleteResult, err := tools.delete(context.Background(), newToolRequest("delete_project", map[string]any{
		"project_id": project.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, deleteResult.IsError)
	require.Equal(t, "Project Deleted!", textFromResult(t, deleteResult))
}
