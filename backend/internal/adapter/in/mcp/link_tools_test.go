package mcp

import (
	"context"
	"testing"

	"devaulty-backend/internal/dto"

	"github.com/stretchr/testify/require"
)

func TestLinkTools_CreateListGetUpdateDelete(t *testing.T) {
	app := setupMCPTestApp(t)
	project := app.createProject(t, "Link project")
	tools := NewLinkTools(app.linkUseCase)

	createResult, err := tools.create(context.Background(), newToolRequest("create_link", map[string]any{
		"projectID":   project.ID.String(),
		"title":       "Documentation",
		"url":         "https://example.com/docs",
		"description": "Official docs",
	}))
	require.NoError(t, err)
	require.False(t, createResult.IsError)

	created, ok := createResult.StructuredContent.(*dto.LinkView)
	require.True(t, ok)
	require.Equal(t, "Documentation", created.Title)

	listResult, err := tools.list(context.Background(), newToolRequest("list_links", map[string]any{
		"projectID":  project.ID.String(),
		"pageNumber": 1,
		"pageSize":   10,
	}))
	require.NoError(t, err)
	require.False(t, listResult.IsError)
	require.NotNil(t, listResult.StructuredContent)

	getResult, err := tools.get(context.Background(), newToolRequest("get_link", map[string]any{
		"projectID": project.ID.String(),
		"id":        created.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, getResult.IsError)

	fetched, ok := getResult.StructuredContent.(*dto.LinkView)
	require.True(t, ok)
	require.Equal(t, created.ID, fetched.ID)

	updateResult, err := tools.update(context.Background(), newToolRequest("update_link", map[string]any{
		"projectID":   project.ID.String(),
		"id":          created.ID.String(),
		"title":       "Updated docs",
		"url":         "https://example.com/updated-docs",
		"description": "Updated official docs",
	}))
	require.NoError(t, err)
	require.False(t, updateResult.IsError)

	updated, ok := updateResult.StructuredContent.(*dto.LinkView)
	require.True(t, ok)
	require.Equal(t, "Updated docs", updated.Title)

	deleteResult, err := tools.delete(context.Background(), newToolRequest("delete_link", map[string]any{
		"projectID": project.ID.String(),
		"id":        created.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, deleteResult.IsError)
	require.Equal(t, "Link Deleted!", textFromResult(t, deleteResult))

	invalidResult, err := tools.get(context.Background(), newToolRequest("get_link", map[string]any{
		"projectID": project.ID.String(),
		"id":        "not-a-uuid",
	}))
	require.NoError(t, err)
	require.True(t, invalidResult.IsError)
	require.Contains(t, textFromResult(t, invalidResult), "invalid UUID format")
}
