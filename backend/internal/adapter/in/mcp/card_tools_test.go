package mcp

import (
	"context"
	"testing"
	"time"

	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/dto"

	"github.com/stretchr/testify/require"
)

func TestCardTools_CreateListGetUpdateMoveDelete(t *testing.T) {
	app := setupMCPTestApp(t)
	project := app.createProject(t, "Card project")
	board := app.createBoard(t, project.ID, "Delivery board")
	columnA := app.createBoardColumn(t, project.ID, board.ID, "Backlog")
	columnB := app.createBoardColumn(t, project.ID, board.ID, "Doing")
	snippet := app.createSnippet(t, project.ID, "API auth")
	tools := NewCardTools(app.cardUseCase)

	createResult, err := tools.create(context.Background(), newToolRequest("create_card", map[string]any{
		"projectID":   project.ID.String(),
		"boardID":     board.ID.String(),
		"columnID":    columnA.ID.String(),
		"title":       "Build login endpoint",
		"description": "Use the auth snippet `@[Auth helper](item:SNIPPET:" + snippet.ID.String() + ")` in the implementation.",
		"priority":    string(model.CardPriorityHigh),
		"dueDate":     time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"linkedItems": []map[string]any{{"itemType": string(model.ItemTypeSnippet), "itemId": snippet.ID.String()}},
	}))
	require.NoError(t, err)
	require.False(t, createResult.IsError)

	created, ok := createResult.StructuredContent.(*dto.CardView)
	require.True(t, ok)
	require.Equal(t, "Build login endpoint", created.Title)
	require.Len(t, created.LinkedItems, 1)

	listResult, err := tools.list(context.Background(), newToolRequest("list_cards", map[string]any{
		"projectID": project.ID.String(),
		"boardID":   board.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, listResult.IsError)
	require.NotNil(t, listResult.StructuredContent)

	getResult, err := tools.get(context.Background(), newToolRequest("get_card", map[string]any{
		"projectID": project.ID.String(),
		"boardID":   board.ID.String(),
		"id":        created.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, getResult.IsError)

	fetched, ok := getResult.StructuredContent.(*dto.CardView)
	require.True(t, ok)
	require.Equal(t, created.ID, fetched.ID)

	updateResult, err := tools.update(context.Background(), newToolRequest("update_card", map[string]any{
		"id":          created.ID.String(),
		"projectID":   project.ID.String(),
		"boardID":     board.ID.String(),
		"title":       "Build login and refresh endpoint",
		"description": "Updated implementation notes.",
		"priority":    string(model.CardPriorityExtremelyHigh),
	}))
	require.NoError(t, err)
	require.False(t, updateResult.IsError)

	updated, ok := updateResult.StructuredContent.(*dto.CardView)
	require.True(t, ok)
	require.Equal(t, "Build login and refresh endpoint", updated.Title)

	moveResult, err := tools.move(context.Background(), newToolRequest("move_card", map[string]any{
		"id":             created.ID.String(),
		"projectID":      project.ID.String(),
		"boardID":        board.ID.String(),
		"targetColumnID": columnB.ID.String(),
		"position":       0,
	}))
	require.NoError(t, err)
	require.False(t, moveResult.IsError)
	require.Equal(t, "Card Moved!", textFromResult(t, moveResult))

	deleteResult, err := tools.delete(context.Background(), newToolRequest("delete_card", map[string]any{
		"projectID": project.ID.String(),
		"boardID":   board.ID.String(),
		"id":        created.ID.String(),
	}))
	require.NoError(t, err)
	require.False(t, deleteResult.IsError)
	require.Equal(t, "Card Deleted!", textFromResult(t, deleteResult))

	invalidResult, err := tools.get(context.Background(), newToolRequest("get_card", map[string]any{
		"projectID": project.ID.String(),
		"boardID":   board.ID.String(),
		"id":        "not-a-uuid",
	}))
	require.NoError(t, err)
	require.True(t, invalidResult.IsError)
	require.Contains(t, textFromResult(t, invalidResult), "invalid UUID format")
}
