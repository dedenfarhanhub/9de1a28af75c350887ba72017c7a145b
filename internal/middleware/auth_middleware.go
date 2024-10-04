package middleware

import (
	"net/http"

	"github.com/dedenfarhanhub/blog-service/config"
	"github.com/dedenfarhanhub/blog-service/internal/helpers"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware check user
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is missing"})
			c.Abort()
			return
		}

		cfg := config.LoadConfig()
		claims, err := helpers.ValidateToken(tokenString, []byte(cfg.JWTSecret))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Set user ID and email to context
		c.Set("userID", claims.ID)
		c.Set("userEmail", claims.Email)

		c.Next()
	}
}
