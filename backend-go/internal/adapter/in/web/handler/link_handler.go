package handler

import (
	"devaulty-backend/internal/adapter/in/web/common"
	"devaulty-backend/internal/dto"
	"devaulty-backend/internal/usecase"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LinkHandler struct {
	linkUseCase *usecase.LinkUseCase
}

func NewLinkHandler(linkUseCase *usecase.LinkUseCase) *LinkHandler {
	return &LinkHandler{linkUseCase: linkUseCase}
}

func (h *LinkHandler) Create(c *gin.Context) {
	var cmd dto.CreateLinkCommand
	err := c.ShouldBindJSON(&cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ProjectID = projectID

	link, err := h.linkUseCase.Create(c.Request.Context(), cmd)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Printf("[LinkHandler.Create] %v", err)
		return
	}
	location := fmt.Sprintf("%s/%s", c.Request.URL.Path, link.ID)
	c.Header("Location", location)
	c.JSON(http.StatusCreated, link)
}

func (h *LinkHandler) GetAll(c *gin.Context) {
	var query common.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pagedLinks, err := h.linkUseCase.GetAllByProjectID(c.Request.Context(), projectID, query.PageNumber, query.PageSize)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Printf("[LinkHandler.GetAll] %v", err)
		return
	}
	c.JSON(http.StatusOK, pagedLinks)
}

func (h *LinkHandler) Get(c *gin.Context) {
	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	linkID, err := common.ExtractUUIDParam(c, "link_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link, err := h.linkUseCase.GetByID(c.Request.Context(), projectID, linkID)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) || errors.Is(err, usecase.ErrLinkNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Printf("[LinkHandler.Get] %v", err)
		return
	}
	c.JSON(http.StatusOK, link)
}

func (h *LinkHandler) Update(c *gin.Context) {
	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	linkID, err := common.ExtractUUIDParam(c, "link_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cmd dto.UpdateLinkCommand
	err = c.ShouldBindJSON(&cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ProjectID = projectID
	cmd.ID = linkID

	link, err := h.linkUseCase.Update(c.Request.Context(), cmd)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) || errors.Is(err, usecase.ErrLinkNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Printf("[LinkHandler.Update] %v", err)
		return
	}
	c.JSON(http.StatusOK, link)
}

func (h *LinkHandler) Delete(c *gin.Context) {
	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := common.ExtractUUIDParam(c, "link_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.linkUseCase.Delete(c.Request.Context(), projectID, id)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) || errors.Is(err, usecase.ErrLinkNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[LinkHandler.Delete] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
