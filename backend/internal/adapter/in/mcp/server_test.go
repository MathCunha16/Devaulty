package mcp

import (
	"testing"

	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/require"
)

func TestMCPServerRegistersAllTools(t *testing.T) {
	app := setupMCPTestApp(t)
	server := mcpserver.NewMCPServer("Devaulty MCP Server", "test-version", mcpserver.WithRecovery())

	NewProjectTools(app.projectUseCase).Register(server, Options{})
	NewSnippetTools(app.snippetUseCase).Register(server, Options{})
	NewProblemTools(app.problemUseCase).Register(server, Options{})
	NewLinkTools(app.linkUseCase).Register(server, Options{})
	NewNoteTools(app.noteUseCase).Register(server, Options{})
	NewBoardTools(app.boardUseCase).Register(server, Options{})
	NewBoardColumnTools(app.boardColumnUseCase).Register(server, Options{})
	NewCardTools(app.cardUseCase).Register(server, Options{})
	NewTagTools(app.tagUseCase).Register(server, Options{})
	NewItemTagTools(app.itemTagUseCase).Register(server, Options{})

	tools := server.ListTools()
	require.NotEmpty(t, tools)

	for _, expected := range []string{
		"list_projects", "get_project", "create_project", "update_project", "archive_project", "unarchive_project", "delete_project",
		"list_snippets", "get_snippet", "create_snippet", "update_snippet", "delete_snippet",
		"list_problems", "get_problem", "create_problem", "update_problem", "update_problem_status", "delete_problem",
		"list_links", "get_link", "create_link", "update_link", "delete_link",
		"list_notes", "get_note", "create_note", "update_note", "delete_note",
		"list_boards", "get_board", "get_default_board", "create_board", "update_board", "delete_board",
		"list_board_columns", "get_board_column", "create_board_column", "update_board_column", "reorder_board_columns", "delete_board_column",
		"list_cards", "get_card", "create_card", "update_card", "move_card", "delete_card",
		"list_tags", "get_tag", "search_tags_by_name", "create_tag", "update_tag", "delete_tag",
		"associate_tag_to_item", "dissociate_tag_from_item",
	} {
		_, ok := tools[expected]
		require.Truef(t, ok, "expected tool %q to be registered", expected)
	}
}

func TestMCPServerReadOnlySkipsWriteTools(t *testing.T) {
	app := setupMCPTestApp(t)
	server := mcpserver.NewMCPServer("Devaulty MCP Server", "test-version", mcpserver.WithRecovery())

	NewProjectTools(app.projectUseCase).Register(server, Options{ReadOnly: true})
	NewSnippetTools(app.snippetUseCase).Register(server, Options{ReadOnly: true})
	NewProblemTools(app.problemUseCase).Register(server, Options{ReadOnly: true})
	NewLinkTools(app.linkUseCase).Register(server, Options{ReadOnly: true})
	NewNoteTools(app.noteUseCase).Register(server, Options{ReadOnly: true})
	NewBoardTools(app.boardUseCase).Register(server, Options{ReadOnly: true})
	NewBoardColumnTools(app.boardColumnUseCase).Register(server, Options{ReadOnly: true})
	NewCardTools(app.cardUseCase).Register(server, Options{ReadOnly: true})
	NewTagTools(app.tagUseCase).Register(server, Options{ReadOnly: true})
	NewItemTagTools(app.itemTagUseCase).Register(server, Options{ReadOnly: true})

	tools := server.ListTools()
	require.NotContains(t, tools, "create_project")
	require.NotContains(t, tools, "update_project")
	require.NotContains(t, tools, "archive_project")
	require.NotContains(t, tools, "delete_project")
	require.NotContains(t, tools, "create_tag")
	require.NotContains(t, tools, "update_tag")
	require.NotContains(t, tools, "delete_tag")
	require.NotContains(t, tools, "associate_tag_to_item")
	require.NotContains(t, tools, "dissociate_tag_from_item")

	require.Contains(t, tools, "list_projects")
	require.Contains(t, tools, "get_project")
	require.Contains(t, tools, "search_tags_by_name")
}

func TestMCPServerDeleteFlagsRemoveDeleteTools(t *testing.T) {
	app := setupMCPTestApp(t)
	server := mcpserver.NewMCPServer("Devaulty MCP Server", "test-version", mcpserver.WithRecovery())

	NewProjectTools(app.projectUseCase).Register(server, Options{DisableDelete: true})
	NewSnippetTools(app.snippetUseCase).Register(server, Options{DisableDelete: true})
	NewProblemTools(app.problemUseCase).Register(server, Options{DisableDelete: true})
	NewLinkTools(app.linkUseCase).Register(server, Options{DisableDelete: true})
	NewNoteTools(app.noteUseCase).Register(server, Options{DisableDelete: true})
	NewBoardTools(app.boardUseCase).Register(server, Options{DisableDelete: true})
	NewBoardColumnTools(app.boardColumnUseCase).Register(server, Options{DisableDelete: true})
	NewCardTools(app.cardUseCase).Register(server, Options{DisableDelete: true})
	NewTagTools(app.tagUseCase).Register(server, Options{DisableDelete: true})
	NewItemTagTools(app.itemTagUseCase).Register(server, Options{DisableDelete: true})

	tools := server.ListTools()
	require.NotContains(t, tools, "delete_project")
	require.NotContains(t, tools, "delete_snippet")
	require.NotContains(t, tools, "delete_problem")
	require.NotContains(t, tools, "delete_link")
	require.NotContains(t, tools, "delete_note")
	require.NotContains(t, tools, "delete_board")
	require.NotContains(t, tools, "delete_board_column")
	require.NotContains(t, tools, "delete_card")
	require.NotContains(t, tools, "delete_tag")

	require.Contains(t, tools, "create_project")
	require.Contains(t, tools, "create_snippet")
}
