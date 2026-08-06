package handler

import (
	"devaulty-backend/internal/adapter/in/web/common"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/usecase"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ItemTagHandler struct {
	itemTagUseCase *usecase.ItemTagUseCase
	tagUseCase     *usecase.TagUseCase
	projectUseCase *usecase.ProjectUseCase
}

func NewItemTagHandler(itemTagUseCase *usecase.ItemTagUseCase, tagUseCase *usecase.TagUseCase, projectUseCase *usecase.ProjectUseCase) *ItemTagHandler {
	return &ItemTagHandler{itemTagUseCase: itemTagUseCase, tagUseCase: tagUseCase, projectUseCase: projectUseCase}
}

func (h *ItemTagHandler) Associate(c *gin.Context) {
	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	itemID, err := common.ExtractUUIDParam(c, "item_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tagID, err := common.ExtractUUIDParam(c, "tag_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	itemType := model.ItemType(strings.ToUpper(c.Param("item_type")))

	err = h.itemTagUseCase.AssociateTagToItem(c.Request.Context(), projectID, itemType, itemID, tagID)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) || errors.Is(err, usecase.ErrTagNotFound) || errors.Is(err, usecase.ErrItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, usecase.ErrUnsupportedItemType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[ItemTagHandler.Associate] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ItemTagHandler) Disassociate(c *gin.Context) {
	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	itemID, err := common.ExtractUUIDParam(c, "item_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tagID, err := common.ExtractUUIDParam(c, "tag_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	itemType := model.ItemType(strings.ToUpper(c.Param("item_type")))

	err = h.itemTagUseCase.DisassociateTagFromItem(c.Request.Context(), projectID, itemType, itemID, tagID)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) || errors.Is(err, usecase.ErrTagNotFound) || errors.Is(err, usecase.ErrItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, usecase.ErrUnsupportedItemType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[ItemTagHandler.Disassociate] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.Status(http.StatusNoContent)
}
