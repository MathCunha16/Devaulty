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

type NoteHandler struct {
	noteUseCase *usecase.NoteUseCase
}

func NewNoteHandler(noteUseCase *usecase.NoteUseCase) *NoteHandler {
	return &NoteHandler{noteUseCase: noteUseCase}
}

func (h *NoteHandler) Create(c *gin.Context) {
	var cmd dto.CreateNoteCommand
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

	note, err := h.noteUseCase.Create(c.Request.Context(), cmd)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	location := fmt.Sprintf("%s/%s", c.Request.URL.Path, note.ID)
	c.Header("Location", location)
	c.JSON(http.StatusCreated, note)
}

func (h *NoteHandler) GetAll(c *gin.Context) {
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

	paginatedNotes, err := h.noteUseCase.GetAllByProjectID(c.Request.Context(), projectID, query.PageNumber, query.PageSize)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[NoteHandler.GetAll] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, paginatedNotes)
}

func (h *NoteHandler) Get(c *gin.Context) {
	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := common.ExtractUUIDParam(c, "note_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.noteUseCase.GetByID(c.Request.Context(), projectID, id)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) || errors.Is(err, usecase.ErrNoteNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[NoteHandler.Get] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if note == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) Update(c *gin.Context) {
	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := common.ExtractUUIDParam(c, "note_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cmd dto.UpdateNoteCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ProjectID = projectID
	cmd.ID = id

	note, err := h.noteUseCase.Update(c.Request.Context(), cmd)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) || errors.Is(err, usecase.ErrNoteNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[NoteHandler.Update] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) Delete(c *gin.Context) {
	projectID, err := common.ExtractUUIDParam(c, "project_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := common.ExtractUUIDParam(c, "note_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.noteUseCase.Delete(c.Request.Context(), projectID, id)
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) || errors.Is(err, usecase.ErrNoteNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[NoteHandler.Delete] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
