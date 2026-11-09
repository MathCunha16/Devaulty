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

type DrawingHandler struct {
	drawingUseCase *usecase.DrawingUseCase
}

func NewDrawingHandler(drawingUseCase *usecase.DrawingUseCase) *DrawingHandler {
	return &DrawingHandler{
		drawingUseCase: drawingUseCase,
	}
}

func (h *DrawingHandler) Create(c *gin.Context) {
	var cmd dto.CreateDrawingCommand
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

	drawing, err := h.drawingUseCase.Create(c.Request.Context(), cmd)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[DrawingHandler.Create] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	location := fmt.Sprintf("%s/%s", c.Request.URL.Path, drawing.ID)
	c.Header("Location", location)
	c.JSON(http.StatusCreated, drawing)
}

func (h *DrawingHandler) GetAll(c *gin.Context) {
	var query common.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paginatedDrawings, err := h.drawingUseCase.GetAllByProjectID(c.Request.Context(), projectID, query.PageNumber, query.PageSize)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[DrawingHandler.GetAll] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, paginatedDrawings)
}

func (h *DrawingHandler) Get(c *gin.Context) {
	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := common.ExtractUUIDParam(c, "drawing_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	drawing, err := h.drawingUseCase.GetByID(c.Request.Context(), projectID, id)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) || errors.Is(err, usecase.ErrDrawingNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[DrawingHandler.Get] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if drawing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Drawing not found"})
		return
	}

	c.JSON(http.StatusOK, drawing)
}

func (h *DrawingHandler) Update(c *gin.Context) {
	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := common.ExtractUUIDParam(c, "drawing_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cmd dto.UpdateDrawingCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ProjectID = projectID
	cmd.ID = id

	drawing, err := h.drawingUseCase.Update(c.Request.Context(), cmd)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) || errors.Is(err, usecase.ErrDrawingNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[DrawingHandler.Update] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, drawing)
}

func (h *DrawingHandler) Delete(c *gin.Context) {
	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := common.ExtractUUIDParam(c, "drawing_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.drawingUseCase.Delete(c.Request.Context(), projectID, id)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) || errors.Is(err, usecase.ErrDrawingNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[DrawingHandler.Delete] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
