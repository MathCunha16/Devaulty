package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"devaulty-backend/internal/adapter/in/web/common"
	"devaulty-backend/internal/dto"
	"devaulty-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type SnippetHandler struct {
	snippetUseCase *usecase.SnippetUseCase
}

func NewSnippetHandler(snippetUseCase *usecase.SnippetUseCase) *SnippetHandler {
	return &SnippetHandler{snippetUseCase: snippetUseCase}
}

func (h *SnippetHandler) Create(c *gin.Context) {
	var cmd dto.CreateSnippetCommand
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

	snippet, err := h.snippetUseCase.Create(c.Request.Context(), cmd)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[SnippetHandler.Create] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	location := fmt.Sprintf("%s/%s", c.Request.URL.Path, snippet.ID)
	c.Header("Location", location)
	c.JSON(http.StatusCreated, snippet)
}

func (h *SnippetHandler) GetAll(c *gin.Context) {
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

	pagedSnippets, err := h.snippetUseCase.GetAllByProjectID(c.Request.Context(), projectID, query.PageNumber, query.PageSize)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[SnippetHandler.GetAll] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, pagedSnippets)
}

func (h *SnippetHandler) Get(c *gin.Context) {
	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := common.ExtractUUIDParam(c, "snippet_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	snippet, err := h.snippetUseCase.GetByID(c.Request.Context(), projectID, id)
	if err != nil {
		if errors.Is(err, usecase.ErrSnippetNotFound) || errors.Is(err, usecase.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[SnippetHandler.Get] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if snippet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Snippet not found"})
		return
	}

	c.JSON(http.StatusOK, snippet)
}

func (h *SnippetHandler) Update(c *gin.Context) {
	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := common.ExtractUUIDParam(c, "snippet_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cmd dto.UpdateSnippetCommand
	err = c.ShouldBindJSON(&cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ID = id
	cmd.ProjectID = projectID

	snippet, err := h.snippetUseCase.Update(c.Request.Context(), cmd)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) || errors.Is(err, usecase.ErrSnippetNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[SnippetHandler.Update] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, snippet)
}

func (h *SnippetHandler) Delete(c *gin.Context) {
	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := common.ExtractUUIDParam(c, "snippet_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.snippetUseCase.Delete(c.Request.Context(), projectID, id)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) || errors.Is(err, usecase.ErrSnippetNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[SnippetHandler.Delete] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
