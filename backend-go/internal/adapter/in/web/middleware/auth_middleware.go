package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(apiToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.GetHeader("DEVAULTY_INTERNAL_TOKEN")
		if clientToken != apiToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
