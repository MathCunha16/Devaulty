package mcp

import (
	"context"
	"devaulty-backend/internal/adapter/in/mcp/util"
	"devaulty-backend/internal/dto"
	"devaulty-backend/internal/usecase"
	"log"

	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type NoteTools struct {
	noteUseCase *usecase.NoteUseCase
}

func NewNoteTools(noteUseCase *usecase.NoteUseCase) *NoteTools {
	return &NoteTools{noteUseCase: noteUseCase}
}

func (t *NoteTools) Register(s *server.MCPServer, opts Options) {

	listNoteTool := mcp.NewTool("list_notes",
		mcp.WithDescription("(paginated) (first page at 0) lists all notes in a Devaulty project. Return a summary view of the notes"),
		mcp.WithInputSchema[util.ProjectPaginationQuery](),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations))
	s.AddTool(listNoteTool, t.list)

	getNoteTool := mcp.NewTool("get_note",
		mcp.WithDescription("gets a note by ID (uuid)"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the note")),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations))
	s.AddTool(getNoteTool, t.get)

	if opts.ReadOnly {
		return
	}

	createNoteTool := mcp.NewTool("create_note",
		mcp.WithDescription("creates a new note in a Devaulty project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("title", mcp.Required(), mcp.Description("The title of the note")),
		mcp.WithString("content", mcp.Required(), mcp.Description("The content of the note. Supports Markdown!")),
		mcp.WithToolAnnotation(util.WriteAnnotations),
	)
	s.AddTool(createNoteTool, t.create)

	updateNoteTool := mcp.NewTool("update_note",
		mcp.WithDescription("updates an existing note in a Devaulty project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the note")),
		mcp.WithString("title", mcp.Description("The title of the note")),
		mcp.WithString("content", mcp.Description("The content of the note. Supports Markdown!")),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(updateNoteTool, t.update)

	if opts.DisableDelete {
		return
	}

	deleteNoteTool := mcp.NewTool("delete_note",
		mcp.WithDescription("deletes an existing note in a Devaulty project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the note")),
		mcp.WithToolAnnotation(util.DeleteAnnotations))
	s.AddTool(deleteNoteTool, t.delete)
}

func (t *NoteTools) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.CreateNoteCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	note, err := t.noteUseCase.Create(ctx, cmd)
	if err != nil {
		log.Printf("[NoteTools.create] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(note)
}

func (t *NoteTools) list(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var query util.ProjectPaginationQuery
	if err := req.BindArguments(&query); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	util.ValidateProjectQuery(&query)

	projectID, err := uuid.Parse(query.ProjectID)
	if err != nil {
		return mcp.NewToolResultError("Invalid project UUID: " + err.Error()), nil
	}

	pagedNotes, err := t.noteUseCase.GetAllByProjectID(ctx, projectID, query.PageNumber, query.PageSize)
	if err != nil {
		log.Printf("[NoteTools.list] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(pagedNotes)
}

func (t *NoteTools) get(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	noteID, err := util.ExtractUUID(req, "id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	note, err := t.noteUseCase.GetByID(ctx, projectID, noteID)
	if err != nil {
		log.Printf("[NoteTools.get] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(note)
}

func (t *NoteTools) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.UpdateNoteCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	note, err := t.noteUseCase.Update(ctx, cmd)
	if err != nil {
		log.Printf("[NoteTools.update] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(note)
}

func (t *NoteTools) delete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	noteID, err := util.ExtractUUID(req, "id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = t.noteUseCase.Delete(ctx, projectID, noteID)
	if err != nil {
		log.Printf("[NoteTools.delete] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Note Deleted!"), nil
}
