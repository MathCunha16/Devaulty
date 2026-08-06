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

type TagHandler struct {
	tagUseCase *usecase.TagUseCase
}

func NewTagHandler(tagUseCase *usecase.TagUseCase) *TagHandler {
	return &TagHandler{tagUseCase: tagUseCase}
}

func (h *TagHandler) Create(c *gin.Context) {
	var cmd dto.CreateTagCommand
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

	tag, err := h.tagUseCase.Create(c.Request.Context(), cmd)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, usecase.ErrTagAlreadyExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[TagHandler.Create] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	location := fmt.Sprintf("%s/%s", c.Request.URL.Path, tag.ID)
	c.Header("Location", location)
	c.JSON(http.StatusCreated, tag)
}

func (h *TagHandler) GetAll(c *gin.Context) {
	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tags, err := h.tagUseCase.GetAllByProjectID(c.Request.Context(), projectID)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[TagHandler.GetAll] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

func (h *TagHandler) Get(c *gin.Context) {
	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := common.ExtractUUIDParam(c, "tag_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, err := h.tagUseCase.GetByID(c.Request.Context(), projectID, id)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) || errors.Is(err, usecase.ErrTagNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[TagHandler.Get] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if tag == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (h *TagHandler) SearchByName(c *gin.Context) {
	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name := c.Query("tag_name")
	tags, err := h.tagUseCase.SearchByName(c.Request.Context(), projectID, name)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[TagHandler.SearchByName] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

func (h *TagHandler) Update(c *gin.Context) {
	var cmd dto.UpdateTagCommand
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
	id, err := common.ExtractUUIDParam(c, "tag_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ID = id
	cmd.ProjectID = projectID

	tag, err := h.tagUseCase.Update(c.Request.Context(), cmd)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) || errors.Is(err, usecase.ErrTagNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, usecase.ErrTagAlreadyExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[TagHandler.Update] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (h *TagHandler) Delete(c *gin.Context) {
	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := common.ExtractUUIDParam(c, "tag_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.tagUseCase.Delete(c.Request.Context(), projectID, id)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) || errors.Is(err, usecase.ErrTagNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[TagHandler.Delete] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
