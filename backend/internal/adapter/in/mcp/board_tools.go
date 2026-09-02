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

type BoardTools struct {
	boardUseCase *usecase.BoardUseCase
}

func NewBoardTools(boardUseCase *usecase.BoardUseCase) *BoardTools {
	return &BoardTools{boardUseCase: boardUseCase}
}

func (t *BoardTools) Register(s *server.MCPServer, opts Options) {

	listBoardTool := mcp.NewTool("list_boards",
		mcp.WithDescription("(paginated) (first page at 0) lists all boards in a Devaulty project"),
		mcp.WithInputSchema[util.ProjectPaginationQuery](),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations))
	s.AddTool(listBoardTool, t.list)

	getBoardTool := mcp.NewTool("get_board",
		mcp.WithDescription("gets a board by ID (uuid)"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the board")),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations))
	s.AddTool(getBoardTool, t.get)

	getDefaultBoardTool := mcp.NewTool("get_default_board",
		mcp.WithDescription("gets the default board for a project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations))
	s.AddTool(getDefaultBoardTool, t.getDefault)

	if opts.ReadOnly {
		return
	}

	createBoardTool := mcp.NewTool("create_board",
		mcp.WithDescription("creates a new kanban board in a Devaulty project (with default columns 'Backlog', 'In Progress', 'Review', 'Done')"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("name", mcp.Required(), mcp.Description("The name of the board")),
		mcp.WithString("description", mcp.Description("The description of the board")),
		mcp.WithBoolean("isDefault", mcp.Description("Whether the board is the default board for the project. If this is true, any other default board will be automatically set to false.")),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(createBoardTool, t.create)

	updateBoardTool := mcp.NewTool("update_board",
		mcp.WithDescription("updates an existing board in a Devaulty project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the board")),
		mcp.WithString("name", mcp.Description("The name of the board")),
		mcp.WithString("description", mcp.Description("The description of the board")),
		mcp.WithBoolean("isDefault", mcp.Description("Whether the board is the default board for the project. If this is true, any other default board will be automatically set to false.")),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(updateBoardTool, t.update)

	if opts.DisableDelete {
		return
	}

	deleteBoardTool := mcp.NewTool("delete_board",
		mcp.WithDescription("deletes an existing board and all of its columns and cards. This action is irreversible."),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the board")),
		mcp.WithToolAnnotation(util.DeleteAnnotations))
	s.AddTool(deleteBoardTool, t.delete)
}

func (t *BoardTools) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.CreateBoardCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	board, err := t.boardUseCase.Create(ctx, cmd)
	if err != nil {
		log.Printf("[BoardTools.create] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(board)
}

func (t *BoardTools) list(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var query util.ProjectPaginationQuery
	if err := req.BindArguments(&query); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	util.ValidateProjectQuery(&query)

	projectID, err := uuid.Parse(query.ProjectID)
	if err != nil {
		return mcp.NewToolResultError("Invalid project UUID" + err.Error()), nil
	}

	pagedBoards, err := t.boardUseCase.GetAllByProjectID(ctx, projectID, query.PageNumber, query.PageSize)
	if err != nil {
		log.Printf("[BoardTools.list] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(pagedBoards)
}

func (t *BoardTools) get(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	boardID, err := util.ExtractUUID(req, "id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	board, err := t.boardUseCase.GetByID(ctx, projectID, boardID)
	if err != nil {
		log.Printf("[BoardTools.get] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(board)
}

func (t *BoardTools) getDefault(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	defaultBoard, err := t.boardUseCase.GetDefaultByProjectID(ctx, projectID)
	if err != nil {
		log.Printf("[BoardTools.getDefault] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(defaultBoard)
}

func (t *BoardTools) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.UpdateBoardCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	board, err := t.boardUseCase.Update(ctx, cmd)
	if err != nil {
		log.Printf("[BoardTools.update] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(board)
}

func (t *BoardTools) delete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	boardID, err := util.ExtractUUID(req, "id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = t.boardUseCase.Delete(ctx, projectID, boardID)
	if err != nil {
		log.Printf("[BoardTools.delete] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Board Deleted!"), nil
}
