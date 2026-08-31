package mcp

import (
	"context"
	"testing"

	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/dto"

	"github.com/stretchr/testify/require"
)

func TestSnippetTools_CreateListGetUpdateDelete(t *testing.T) {
	app := setupMCPTestApp(t)
	project := app.createProject(t, "Snippet project")
	tools := NewSnippetTools(app.snippetUseCase)

	createResult, err := tools.create(context.Background(), newToolRequest("create_snippet", map[string]any{
		"projectID":   project.ID.String(),
		"title":       "Bootstrap snippet",
		"description": "A snippet used by the MCP tests",
		"content":     "console.log('hello from mcp')",
		"language":    string(model.SnippetLangJavascript),
		"snippetType": string(model.SnippetTypeCode),
	}))
	require.NoError(t, err)
	require.False(t, createResult.IsError)

	created, ok := createResult.StructuredContent.(*dto.SnippetView)
	require.True(t, ok)
	require.Equal(t, "Bootstrap snippet", created.Title)

	listResult, err := tools.list(context.Background(), newToolRequest("list_snippets", map[string]any{
		"projectID":  project.ID.String(),
		"pageNumber": 1,
		"pageSize":   10,
	}))
	require.NoError(t, err)
	require.False(t, listResult.IsError)
	require.NotNil(t, listResult.StructuredContent)

	getResult, err := tools.get(context.Background(), newToolRequest("get_snippet", map[string]any{
		"projectID": project.ID.String(),
		"id":        created.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, getResult.IsError)

	fetched, ok := getResult.StructuredContent.(*dto.SnippetView)
	require.True(t, ok)
	require.Equal(t, created.ID, fetched.ID)

	updateResult, err := tools.update(context.Background(), newToolRequest("update_snippet", map[string]any{
		"projectID":   project.ID.String(),
		"id":          created.ID.String(),
		"title":       "Updated snippet title",
		"description": "Updated snippet description",
		"content":     "console.log('updated')",
		"language":    string(model.SnippetLangTypescript),
		"snippetType": string(model.SnippetTypeCode),
	}))
	require.NoError(t, err)
	require.False(t, updateResult.IsError)

	updated, ok := updateResult.StructuredContent.(*dto.SnippetView)
	require.True(t, ok)
	require.Equal(t, "Updated snippet title", updated.Title)

	deleteResult, err := tools.delete(context.Background(), newToolRequest("delete_snippet", map[string]any{
		"projectID": project.ID.String(),
		"id":        created.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, deleteResult.IsError)
	require.Equal(t, "Snippet Deleted!", textFromResult(t, deleteResult))

	invalidResult, err := tools.get(context.Background(), newToolRequest("get_snippet", map[string]any{
		"projectID": "not-a-uuid",
		"id":        created.ID.String(),
	}))
	require.NoError(t, err)
	require.True(t, invalidResult.IsError)
	require.Contains(t, textFromResult(t, invalidResult), "invalid UUID format")
}
