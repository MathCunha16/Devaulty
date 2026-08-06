package web

import (
	"devaulty-backend/internal/adapter/in/web/handler"
	"devaulty-backend/internal/adapter/in/web/middleware"
	"net/http"
	"os"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Project *handler.ProjectHandler
	Snippet *handler.SnippetHandler
	Link    *handler.LinkHandler
	Problem *handler.ProblemHandler
	Tag     *handler.TagHandler
	ItemTag *handler.ItemTagHandler
}

func SetupRouter(h *Handlers, apiToken string) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	registerDocsRoutes(r)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	v1 := r.Group("/api/v1")
	{
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(apiToken))
		{
			mapProjectRoutes(protected, h.Project)
			mapSnippetRoutes(protected, h.Snippet)
			mapLinkRoutes(protected, h.Link)
			mapProblemRoutes(protected, h.Problem)
			mapTagRoutes(protected, h.Tag)
			mapItemTagRoutes(protected, h.ItemTag)
		}
	}
	return r
}

func mapProjectRoutes(rg *gin.RouterGroup, h *handler.ProjectHandler) {
	projects := rg.Group("/projects")
	{
		projects.POST("", h.Create)
		projects.GET("", h.GetAll)
		projects.GET("/:project_id", h.Get)
		projects.PATCH("/:project_id", h.Update)
		projects.PATCH("/:project_id/archive", h.Archive)
		projects.PATCH("/:project_id/unarchive", h.Unarchive)
		projects.DELETE("/:project_id", h.Delete)
	}
}

func mapSnippetRoutes(rg *gin.RouterGroup, h *handler.SnippetHandler) {
	snippets := rg.Group("/projects/:project_id/snippets")
	{
		snippets.POST("", h.Create)
		snippets.GET("", h.GetAll)
		snippets.GET("/:snippet_id", h.Get)
		snippets.PATCH("/:snippet_id", h.Update)
		snippets.DELETE("/:snippet_id", h.Delete)
	}
}

func mapLinkRoutes(rg *gin.RouterGroup, h *handler.LinkHandler) {
	links := rg.Group("/projects/:project_id/links")
	{
		links.POST("", h.Create)
		links.GET("", h.GetAll)
		links.GET("/:link_id", h.Get)
		links.PATCH("/:link_id", h.Update)
		links.DELETE("/:link_id", h.Delete)
	}
}

func mapProblemRoutes(rg *gin.RouterGroup, h *handler.ProblemHandler) {
	problems := rg.Group("/projects/:project_id/problems")
	{
		problems.POST("", h.Create)
		problems.GET("", h.GetAll)
		problems.GET("/:problem_id", h.Get)
		problems.PATCH("/:problem_id", h.Update)
		problems.PATCH("/:problem_id/status", h.UpdateStatus)
		problems.DELETE("/:problem_id", h.Delete)
	}
}

func mapTagRoutes(rg *gin.RouterGroup, h *handler.TagHandler) {
	tags := rg.Group("/projects/:project_id/tags")
	{
		tags.POST("", h.Create)
		tags.GET("", h.GetAll)
		tags.GET("/search", h.SearchByName)
		tags.GET("/:tag_id", h.Get)
		tags.PATCH("/:tag_id", h.Update)
		tags.DELETE("/:tag_id", h.Delete)
	}
}

func mapItemTagRoutes(rg *gin.RouterGroup, h *handler.ItemTagHandler) {
	itemTags := rg.Group("/projects/:project_id/items/:item_type/:item_id/tags")
	{
		itemTags.PUT("/:tag_id", h.Associate)
		itemTags.DELETE("/:tag_id", h.Disassociate)
	}
}

func registerDocsRoutes(r *gin.Engine) {
	if os.Getenv("APP_ENV") != "dev" {
		return
	}

	r.StaticFile("/openapi.yaml", "./docs/openapi.yaml")

	r.GET("/docs", func(c *gin.Context) {
		specContent, err := os.ReadFile("./docs/openapi.yaml")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read openapi.yaml: "+err.Error())
			return
		}

		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecContent: string(specContent),
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Devaulty API Documentation",
			},
			DarkMode: true,
		})
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
	})
}
