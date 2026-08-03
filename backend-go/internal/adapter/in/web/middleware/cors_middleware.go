package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var allowedOrigins = []string{
	"tauri://localhost",
	"http://tauri.localhost",
	"http://localhost:1420",
	"http://localhost:5173",
	"http://localhost:8080",
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		isAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if strings.EqualFold(origin, allowedOrigin) {
				isAllowed = true
				break
			}
		}
		if isAllowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, "+
			"X-CSRF-Token, Authorization, accept,origin, Cache-Control, X-Requested-With, DEVAULTY_INTERNAL_TOKEN")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
