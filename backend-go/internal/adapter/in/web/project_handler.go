package web

import (
	"devaulty-backend/internal/usecase"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	projectUseCase *usecase.ProjectUseCase
}

func NewProjectHandler(projectUseCase *usecase.ProjectUseCase) *ProjectHandler {
	return &ProjectHandler{projectUseCase: projectUseCase}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var cmd usecase.CreateProjectCommand
	err := c.ShouldBindJSON(&cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	project, err := h.projectUseCase.Create(c, cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	location := fmt.Sprintf("%s/%s", c.Request.URL.Path, project.ID)
	c.Header("Location", location)
	c.JSON(http.StatusCreated, project)
}

func (h *ProjectHandler) Get(c *gin.Context) {
	var idParam = c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := h.projectUseCase.GetByID(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}
