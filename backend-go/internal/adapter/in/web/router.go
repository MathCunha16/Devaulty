package web

import (
	"devaulty-backend/internal/adapter/in/web/middleware"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Project *ProjectHandler
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

func mapProjectRoutes(rg *gin.RouterGroup, h *ProjectHandler) {
	projects := rg.Group("/projects")
	{
		projects.POST("", h.Create)
		projects.GET("/:id", h.Get)
	}
}
