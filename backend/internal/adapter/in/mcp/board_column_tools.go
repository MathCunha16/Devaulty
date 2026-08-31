package mcp

import (
	"context"
	"devaulty-backend/internal/adapter/in/mcp/util"
	"devaulty-backend/internal/dto"
	"devaulty-backend/internal/usecase"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type BoardColumnTools struct {
	boardColumnUseCase *usecase.BoardColumnUseCase
}

func NewBoardColumnTools(boardColumnUseCase *usecase.BoardColumnUseCase) *BoardColumnTools {
	return &BoardColumnTools{boardColumnUseCase: boardColumnUseCase}
}

func (t *BoardColumnTools) Register(s *server.MCPServer, opts Options) {

	listBoardColumnTool := mcp.NewTool("list_board_columns",
		mcp.WithDescription("(non-paged) lists all board columns of a Devaulty board"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("boardID", mcp.Required(), mcp.Description("The boardID (UUID) of the board")),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations))
	s.AddTool(listBoardColumnTool, t.list)

	getBoardColumnTool := mcp.NewTool("get_board_column",
		mcp.WithDescription("gets a board column by ID (uuid)"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("boardID", mcp.Required(), mcp.Description("The boardID (UUID) of the board")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the board column")),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations))
	s.AddTool(getBoardColumnTool, t.get)

	if opts.ReadOnly {
		return
	}

	createBoardColumnTool := mcp.NewTool("create_board_column",
		mcp.WithDescription("creates a new kanban board column in a Devaulty project existent board"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("boardID", mcp.Required(), mcp.Description("The boardID (UUID) of the board")),
		mcp.WithString("name", mcp.Required(), mcp.Description("The name of the column")),
		mcp.WithInteger("wipLimit", mcp.Description("Optional WIP limit for the column. Omit this field to disable the limit.")),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(createBoardColumnTool, t.create)

	updateBoardColumnTool := mcp.NewTool("update_board_column",
		mcp.WithDescription("updates an existing board column in a Devaulty project"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("boardID", mcp.Required(), mcp.Description("The boardID (UUID) of the board")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the board column")),
		mcp.WithString("name", mcp.Description("The name of the column")),
		mcp.WithInteger("wipLimit", mcp.Description("The WIP limit (uint16) of the column.")),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(updateBoardColumnTool, t.update)

	reorderBoardColumnTool := mcp.NewTool("reorder_board_columns",
		mcp.WithDescription("Reorders the columns of a board"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("boardID", mcp.Required(), mcp.Description("The boardID (UUID) of the board")),
		mcp.WithArray("positions", mcp.Required(), mcp.Description("Ordered list of column UUIDs, in the new desired order"), mcp.WithStringItems()),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(reorderBoardColumnTool, t.reorder)

	if opts.DisableDelete {
		return
	}

	deleteBoardColumnTool := mcp.NewTool("delete_board_column",
		mcp.WithDescription("deletes an existing board column in a Devaulty board. This action is irreversible and will delete all cards in the column."),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("boardID", mcp.Required(), mcp.Description("The boardID (UUID) of the board")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the board column")),
		mcp.WithToolAnnotation(util.DeleteAnnotations))
	s.AddTool(deleteBoardColumnTool, t.delete)
}

func (t *BoardColumnTools) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.CreateBoardColumnCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	bc, err := t.boardColumnUseCase.Create(ctx, cmd)
	if err != nil {
		log.Printf("[BoardColumnTools.create] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(bc)
}

func (t *BoardColumnTools) list(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	boardID, err := util.ExtractUUID(req, "boardID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	bcList, err := t.boardColumnUseCase.GetAllByBoardID(ctx, projectID, boardID)
	if err != nil {
		log.Printf("[BoardColumnTools.list] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(bcList)
}

func (t *BoardColumnTools) get(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	boardID, err := util.ExtractUUID(req, "boardID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	bcID, err := util.ExtractUUID(req, "id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	bc, err := t.boardColumnUseCase.GetByIDAndBoardID(ctx, projectID, boardID, bcID)
	if err != nil {
		log.Printf("[BoardColumnTools.get] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(bc)
}

func (t *BoardColumnTools) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.UpdateBoardColumnCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	bc, err := t.boardColumnUseCase.Update(ctx, cmd)
	if err != nil {
		log.Printf("[BoardColumnTools.update] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(bc)
}

func (t *BoardColumnTools) delete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	boardID, err := util.ExtractUUID(req, "boardID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	bcID, err := util.ExtractUUID(req, "id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = t.boardColumnUseCase.Delete(ctx, projectID, boardID, bcID)
	if err != nil {
		log.Printf("[BoardColumnTools.delete] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Board Column Deleted!"), nil
}

func (t *BoardColumnTools) reorder(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.ReorderBoardColumnsCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	bc, err := t.boardColumnUseCase.Reorder(ctx, cmd)
	if err != nil {
		log.Printf("[BoardColumnTools.reorder] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(bc)
}
