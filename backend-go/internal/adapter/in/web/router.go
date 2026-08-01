package web

import (
	"devaulty-backend/internal/adapter/in/web/handler"
	"devaulty-backend/internal/adapter/in/web/middleware"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Project *handler.ProjectHandler
}

func SetupRouter(h *Handlers, apiToken string) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	v1 := r.Group("/api/v1")
	{
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(apiToken))
		{
			mapProjectRoutes(protected, h.Project)
		}
	}
	return r
}

func mapProjectRoutes(rg *gin.RouterGroup, h *handler.ProjectHandler) {
	projects := rg.Group("/projects")
	{
		projects.POST("", h.Create)
		projects.GET("", h.GetAll)
		projects.GET("/:id", h.Get)
		projects.PATCH("/:id", h.Update)
		projects.PATCH("/:id/archive", h.Archive)
		projects.PATCH("/:id/unarchive", h.Unarchive)
		projects.DELETE("/:id", h.Delete)
	}
}
