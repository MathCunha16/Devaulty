package mcp

import (
	"context"
	"testing"

	"devaulty-backend/internal/adapter/out/persistence"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/dto"
	"devaulty-backend/internal/usecase"

	"github.com/google/uuid"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

type mcpTestApp struct {
	projectUseCase     *usecase.ProjectUseCase
	snippetUseCase     *usecase.SnippetUseCase
	problemUseCase     *usecase.ProblemUseCase
	linkUseCase        *usecase.LinkUseCase
	noteUseCase        *usecase.NoteUseCase
	boardUseCase       *usecase.BoardUseCase
	boardColumnUseCase *usecase.BoardColumnUseCase
	cardUseCase        *usecase.CardUseCase
	tagUseCase         *usecase.TagUseCase
	itemTagUseCase     *usecase.ItemTagUseCase
}

func setupMCPTestApp(t *testing.T) *mcpTestApp {
	t.Helper()

	db, err := persistence.InitDB(":memory:", "../../../../migrations")
	require.NoError(t, err)

	itemTagRepo := persistence.NewItemTagRepository(db)
	projectRepo := persistence.NewProjectRepository(db)
	snippetRepo := persistence.NewSnippetRepository(db)
	linkRepo := persistence.NewLinkRepository(db)
	problemRepo := persistence.NewProblemRepository(db)
	noteRepo := persistence.NewNoteRepository(db)
	tagRepo := persistence.NewTagRepository(db)
	credentialRepo := persistence.NewCredentialRepository(db)
	boardRepo := persistence.NewBoardRepository(db)
	boardColumnRepo := persistence.NewBoardColumnRepository(db)
	cardRepo := persistence.NewCardRepository(db)

	return &mcpTestApp{
		projectUseCase:     usecase.NewProjectUseCase(projectRepo),
		snippetUseCase:     usecase.NewSnippetUseCase(snippetRepo, projectRepo, itemTagRepo),
		problemUseCase:     usecase.NewProblemUseCase(problemRepo, projectRepo, itemTagRepo),
		linkUseCase:        usecase.NewLinkUseCase(linkRepo, projectRepo, itemTagRepo),
		noteUseCase:        usecase.NewNoteUseCase(noteRepo, projectRepo, itemTagRepo),
		boardUseCase:       usecase.NewBoardUseCase(boardRepo, boardColumnRepo, projectRepo, itemTagRepo),
		boardColumnUseCase: usecase.NewBoardColumnUseCase(boardColumnRepo, boardRepo, projectRepo),
		cardUseCase:        usecase.NewCardUseCase(cardRepo, boardRepo, boardColumnRepo, projectRepo, itemTagRepo),
		tagUseCase:         usecase.NewTagUseCase(tagRepo, projectRepo),
		itemTagUseCase:     usecase.NewItemTagUseCase(itemTagRepo, tagRepo, projectRepo, snippetRepo, credentialRepo, linkRepo, problemRepo, noteRepo, boardRepo, cardRepo),
	}
}

func newToolRequest(name string, args map[string]any) mcpgo.CallToolRequest {
	return mcpgo.CallToolRequest{Params: mcpgo.CallToolParams{Name: name, Arguments: args}}
}

func (a *mcpTestApp) createProject(t *testing.T, name string) *model.Project {
	t.Helper()

	description := "Project description"
	color := "#10b981"
	icon := "Folder"
	project, err := a.projectUseCase.Create(context.Background(), dto.CreateProjectCommand{
		Name:        name,
		Description: &description,
		Color:       &color,
		Icon:        &icon,
	})
	require.NoError(t, err)
	return project
}

func (a *mcpTestApp) createTag(t *testing.T, projectID uuid.UUID, name string) *dto.TagView {
	t.Helper()

	color := "#ef4444"
	tag, err := a.tagUseCase.Create(context.Background(), dto.CreateTagCommand{
		ProjectID: projectID,
		Name:      name,
		Color:     &color,
	})
	require.NoError(t, err)
	return tag
}

func (a *mcpTestApp) createBoard(t *testing.T, projectID uuid.UUID, name string) *dto.BoardView {
	t.Helper()

	description := "Board description"
	board, err := a.boardUseCase.Create(context.Background(), dto.CreateBoardCommand{
		ProjectID:   projectID,
		Name:        name,
		Description: &description,
		IsDefault:   false,
	})
	require.NoError(t, err)
	return board
}

func (a *mcpTestApp) createBoardColumn(t *testing.T, projectID, boardID uuid.UUID, name string) *dto.BoardColumnView {
	t.Helper()

	wipLimit := uint16(3)
	column, err := a.boardColumnUseCase.Create(context.Background(), dto.CreateBoardColumnCommand{
		ProjectID: projectID,
		BoardID:   boardID,
		Name:      name,
		WipLimit:  &wipLimit,
	})
	require.NoError(t, err)
	return column
}

func (a *mcpTestApp) createSnippet(t *testing.T, projectID uuid.UUID, title string) *dto.SnippetView {
	t.Helper()

	description := "Snippet description"
	snippet, err := a.snippetUseCase.Create(context.Background(), dto.CreateSnippetCommand{
		ProjectID:   projectID,
		Title:       title,
		Description: &description,
		Content:     "console.log('hello')",
		Language:    model.SnippetLangJavascript,
		SnippetType: model.SnippetTypeCode,
	})
	require.NoError(t, err)
	return snippet
}

func (a *mcpTestApp) createNote(t *testing.T, projectID uuid.UUID, title, content string) *dto.NoteView {
	t.Helper()

	note, err := a.noteUseCase.Create(context.Background(), dto.CreateNoteCommand{
		ProjectID: projectID,
		Title:     title,
		Content:   content,
	})
	require.NoError(t, err)
	return note
}

func (a *mcpTestApp) createLink(t *testing.T, projectID uuid.UUID, title, url string) *dto.LinkView {
	t.Helper()

	description := "Link description"
	link, err := a.linkUseCase.Create(context.Background(), dto.CreateLinkCommand{
		ProjectID:   projectID,
		Title:       title,
		URL:         url,
		Description: &description,
	})
	require.NoError(t, err)
	return link
}

func (a *mcpTestApp) createProblem(t *testing.T, projectID uuid.UUID, title string) *dto.ProblemView {
	t.Helper()

	solution := "Check the environment"
	problem, err := a.problemUseCase.Create(context.Background(), dto.CreateProblemCommand{
		ProjectID:        projectID,
		Title:            title,
		ErrorDescription: "Stack trace example",
		Solution:         &solution,
		Status:           model.ProblemStatusOpen,
		Severity:         model.ProblemSeverityHigh,
	})
	require.NoError(t, err)
	return problem
}

func (a *mcpTestApp) createCard(t *testing.T, projectID, boardID, columnID uuid.UUID, title string) *dto.CardView {
	t.Helper()

	description := "Card description"
	priority := model.CardPriorityMedium
	card, err := a.cardUseCase.Create(context.Background(), dto.CreateCardCommand{
		ProjectID:   projectID,
		BoardID:     boardID,
		ColumnID:    columnID,
		Title:       title,
		Description: &description,
		Priority:    &priority,
	})
	require.NoError(t, err)
	return card
}

func textFromResult(t *testing.T, result *mcpgo.CallToolResult) string {
	t.Helper()
	require.NotNil(t, result)
	require.NotEmpty(t, result.Content)
	textContent, ok := result.Content[0].(mcpgo.TextContent)
	require.True(t, ok, "expected text content")
	return textContent.Text
}
