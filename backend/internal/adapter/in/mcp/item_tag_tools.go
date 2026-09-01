package mcp

import (
	"context"
	"devaulty-backend/internal/adapter/in/mcp/util"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/usecase"
	"fmt"
	"log"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var supportedItemTagTypes = []string{
	string(model.ItemTypeSnippet),
	string(model.ItemTypeNote),
	string(model.ItemTypeLink),
	string(model.ItemTypeProblem),
	string(model.ItemTypeBoard),
	string(model.ItemTypeCard),
}

type ItemTagTools struct {
	itemTagUseCase *usecase.ItemTagUseCase
}

func NewItemTagTools(itemTagUseCase *usecase.ItemTagUseCase) *ItemTagTools {
	return &ItemTagTools{itemTagUseCase: itemTagUseCase}
}

func validateSupportedItemType(itemType model.ItemType) error {
	switch itemType {
	case model.ItemTypeSnippet, model.ItemTypeNote, model.ItemTypeLink, model.ItemTypeProblem, model.ItemTypeBoard, model.ItemTypeCard:
		return nil
	default:
		return fmt.Errorf("itemType %q is not supported; supported values are: %s", itemType, strings.Join(supportedItemTagTypes, ", "))
	}
}

func (t *ItemTagTools) Register(s *server.MCPServer, opts Options) {
	if opts.ReadOnly {
		return
	}

	associateTagToItemTool := mcp.NewTool("associate_tag_to_item",
		mcp.WithDescription("Associates a tag to a supported project item (snippet, note, problem, link, card, board)"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("itemID", mcp.Required(), mcp.Description("The itemID (UUID) of the item to tag")),
		mcp.WithString("tagID", mcp.Required(), mcp.Description("The tagID (UUID) of the tag to associate")),
		mcp.WithString("itemType", mcp.Required(), mcp.Description("The type of the item to tag. Supported values: SNIPPET, NOTE, PROBLEM, LINK, BOARD, CARD."), mcp.Enum(supportedItemTagTypes...)),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(associateTagToItemTool, t.associate)

	dissociateTagFromItemTool := mcp.NewTool("dissociate_tag_from_item",
		mcp.WithDescription("Removes a tag association from a supported project item"),
		mcp.WithString("projectID", mcp.Required(), mcp.Description("The projectID (UUID) of the project")),
		mcp.WithString("itemID", mcp.Required(), mcp.Description("The itemID (UUID) of the item to untag")),
		mcp.WithString("tagID", mcp.Required(), mcp.Description("The tagID (UUID) of the tag to remove")),
		mcp.WithString("itemType", mcp.Required(), mcp.Description("The type of the item to untag. Supported values: SNIPPET, NOTE, PROBLEM, LINK, BOARD, CARD."), mcp.Enum(supportedItemTagTypes...)),
		mcp.WithToolAnnotation(util.WriteAnnotations))
	s.AddTool(dissociateTagFromItemTool, t.disassociate)
}

func (t *ItemTagTools) associate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	itemID, err := util.ExtractUUID(req, "itemID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	tagID, err := util.ExtractUUID(req, "tagID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	itemTypeStr, err := req.RequireString("itemType")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	itemType := model.ItemType(strings.ToUpper(itemTypeStr))
	if err := validateSupportedItemType(itemType); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = t.itemTagUseCase.AssociateTagToItem(ctx, projectID, itemType, itemID, tagID)
	if err != nil {
		log.Printf("[ItemTagTools.associate] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Item Tag Associated!"), nil
}

func (t *ItemTagTools) disassociate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := util.ExtractUUID(req, "projectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	itemID, err := util.ExtractUUID(req, "itemID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	tagID, err := util.ExtractUUID(req, "tagID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	itemTypeStr, err := req.RequireString("itemType")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	itemType := model.ItemType(strings.ToUpper(itemTypeStr))
	if err := validateSupportedItemType(itemType); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = t.itemTagUseCase.DisassociateTagFromItem(ctx, projectID, itemType, itemID, tagID)
	if err != nil {
		log.Printf("[ItemTagTools.dissociate] %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Item Tag Dissociated!"), nil
}
