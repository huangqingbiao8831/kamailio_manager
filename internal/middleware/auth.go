package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"your-project/internal/logger"
)

// AuthRequired checks for a specific API Key in the header
func AuthRequired(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-API-Key")

		if token == "" || token != apiKey {
			logger.Warn("Unauthorized access attempt",
				zap.String("ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path))

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid API Key"})
			c.Abort() // cancel next Handler invoke
			return
		}

		c.Next() // continue next Handler invoke
	}
}
