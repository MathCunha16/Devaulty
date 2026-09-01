package mcp

import (
	"context"
	"devaulty-backend/internal/adapter/in/mcp/util"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/dto"
	"devaulty-backend/internal/usecase"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type CardTools struct {
	cardUseCase *usecase.CardUseCase
}

func NewCardTools(cardUseCase *usecase.CardUseCase) *CardTools {
	return &CardTools{cardUseCase: cardUseCase}
}

func (t *CardTools) Register(s *server.MCPServer, opts Options) {

	listCardTool := mcp.NewTool("list_cards",
		mcp.WithDescription("(non-paged) lists all cards in a Devaulty board"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("boardID", mcp.Required(), mcp.Description("The boardID (UUID) of the board")),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations))
	s.AddTool(listCardTool, t.list)

	getCardTool := mcp.NewTool("get_card",
		mcp.WithDescription("gets a card by ID (uuid)"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("boardID", mcp.Required(), mcp.Description("The boardID (UUID) of the board")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the card")),
		mcp.WithToolAnnotation(util.ReadOnlyAnnotations))
	s.AddTool(getCardTool, t.get)

	if opts.ReadOnly {
		return
	}

	createCardTool := mcp.NewTool("create_card",
		mcp.WithDescription("Creates a new card in a Devaulty board. In the card description, you can reference any linked item using Markdown mention syntax: `@[Title](item:TYPE:UUID)`. First add the item to `linkedItems`, then use the same `itemType` and `itemId` in the markdown mention so the frontend renders it as a clickable mention pill."),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("boardID", mcp.Required(), mcp.Description("The boardID (UUID) of the board")),
		mcp.WithString("columnID", mcp.Required(), mcp.Description("The columnID (UUID) of the column")),
		mcp.WithString("title", mcp.Required(), mcp.Description("The title of the card")),
		mcp.WithString("description", mcp.Description("The description of the card in Markdown. If you want to reference a linked item inside the description, attach it first in `linkedItems` and then use the mention format `@[Item Title](item:TYPE:UUID)`, e.g. `@[Login API](item:SNIPPET:123e4567-e89b-12d3-a456-426614174000)` for a snippet or `@[API bug](item:PROBLEM:...)` for a problem. Only items added to `linkedItems` can be mentioned here.")),
		mcp.WithString("priority", mcp.Description("The priority of the card"), mcp.Enum(util.ToStrings(model.CardPriorities)...)),
		mcp.WithString("dueDate", mcp.Description("The Due date in RFC3339/ISO 8601 format, e.g. 2026-08-31T23:59:59Z. Omit if there's no due date.")),
		mcp.WithArray("linkedItems",
			mcp.Description("Optional list of items (snippets, notes, problems, links, etc.) to attach to this card. These are the only items that can be referenced from the card description with `@[Title](item:TYPE:UUID)`. Add the item here before referencing it in the markdown."),
			mcp.Items(map[string]any{
				"type": "object",
				"properties": map[string]any{
					"itemType": map[string]any{
						"type":        "string",
						"enum":        util.ToStrings(model.ItemTypes),
						"description": "Type of the linked item. Use the same value later in the markdown mention, for example `SNIPPET`, `NOTE`, `PROBLEM`, `LINK`, or `BOARD`.",
					},
					"itemId": map[string]any{
						"type":        "string",
						"description": "UUID of the linked item. This is the ID used in the mention syntax `@[Title](item:TYPE:UUID)`.",
					},
				},
				"required": []string{"itemType", "itemId"},
			}),
		),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(createCardTool, t.create)

	updateCardTool := mcp.NewTool("update_card",
		mcp.WithDescription("updates a existent card in a Devaulty board. In the card description, you can reference any linked item using Markdown mention syntax: `@[Title](item:TYPE:UUID)`. First add the item to `linkedItems`, then use the same `itemType` and `itemId` in the markdown mention so the frontend renders it as a clickable mention pill."),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the card")),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("boardID", mcp.Required(), mcp.Description("The boardID (UUID) of the board")),
		mcp.WithString("title", mcp.Description("The title of the card")),
		mcp.WithString("description", mcp.Description("The description of the card in Markdown. If you want to reference a linked item inside the description, attach it first in `linkedItems` and then use the mention format `@[Item Title](item:TYPE:UUID)`, e.g. `@[Login API](item:SNIPPET:123e4567-e89b-12d3-a456-426614174000)` for a snippet or `@[API bug](item:PROBLEM:...)` for a problem. Only items added to `linkedItems` can be mentioned here.")),
		mcp.WithString("priority", mcp.Description("The priority of the card"), mcp.Enum(util.ToStrings(model.CardPriorities)...)),
		mcp.WithString("dueDate", mcp.Description("The Due date in RFC3339/ISO 8601 format, e.g. 2026-08-31T23:59:59Z. Omit if there's no due date.")),
		mcp.WithArray("linkedItems",
			mcp.Description("Optional list of items (snippets, notes, problems, links, etc.) to attach to this card. These are the only items that can be referenced from the card description with `@[Title](item:TYPE:UUID)`. Add the item here before referencing it in the markdown."),
			mcp.Items(map[string]any{
				"type": "object",
				"properties": map[string]any{
					"itemType": map[string]any{
						"type":        "string",
						"enum":        util.ToStrings(model.ItemTypes),
						"description": "Type of the linked item. Use the same value later in the markdown mention, for example `SNIPPET`, `NOTE`, `PROBLEM`, `LINK`, or `BOARD`.",
					},
					"itemId": map[string]any{
						"type":        "string",
						"description": "UUID of the linked item. This is the ID used in the mention syntax `@[Title](item:TYPE:UUID)`.",
					},
				},
				"required": []string{"itemType", "itemId"},
			}),
		),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(updateCardTool, t.update)

	moveCardTool := mcp.NewTool("move_card",
		mcp.WithDescription("moves a card to a new position or a new column in a Devaulty board"),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the card")),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("boardID", mcp.Required(), mcp.Description("The boardID (UUID) of the board")),
		mcp.WithString("targetColumnID", mcp.Required(), mcp.Description("The target columnID (UUID) of the column. Use the same columnID as the current column if you want to move the card to a new position within the same column.")),
		mcp.WithInteger("position", mcp.Description("The position (uint16)(0-based) of the card in the target column. Use 0 if you want to move the card to the beginning of the column."), mcp.Required()),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(moveCardTool, t.move)

	if opts.DisableDelete {
		return
	}

	deleteCardTool := mcp.NewTool("delete_card",
		mcp.WithDescription("deletes an existing card in a Devaulty board. This is irreversible."),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("boardID", mcp.Required(), mcp.Description("The boardID (UUID) of the board")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id (UUID) of the card")),
		mcp.WithToolAnnotation(util.DeleteAnnotations))
	s.AddTool(deleteCardTool, t.delete)

}

func (t *CardTools) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.CreateCardCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	card, err := t.cardUseCase.Create(ctx, cmd)
	if err != nil {
		log.Printf("[CardTools.create] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(card)
}

func (t *CardTools) list(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	boardID, err := util.ExtractUUID(req, "boardID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	cards, err := t.cardUseCase.GetAllByBoardID(ctx, projectID, boardID)
	if err != nil {
		log.Printf("[CardTools.list] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(cards)
}

func (t *CardTools) get(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	boardID, err := util.ExtractUUID(req, "boardID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	cardID, err := util.ExtractUUID(req, "id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	card, err := t.cardUseCase.GetByID(ctx, projectID, boardID, cardID)
	if err != nil {
		log.Printf("[CardTools.get] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(card)
}

func (t *CardTools) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.UpdateCardCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	card, err := t.cardUseCase.Update(ctx, cmd)
	if err != nil {
		log.Printf("[CardTools.update] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(card)
}

func (t *CardTools) delete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	boardID, err := util.ExtractUUID(req, "boardID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	cardID, err := util.ExtractUUID(req, "id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = t.cardUseCase.Delete(ctx, projectID, boardID, cardID)
	if err != nil {
		log.Printf("[CardTools.delete] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Card Deleted!"), nil
}

func (t *CardTools) move(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd dto.MoveCardCommand
	if err := req.BindArguments(&cmd); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err := t.cardUseCase.Move(ctx, cmd)
	if err != nil {
		log.Printf("[CardTools.move] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Card Moved!"), nil
}
