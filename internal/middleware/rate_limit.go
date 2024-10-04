package middleware

import (
	"github.com/dedenfarhanhub/blog-service/internal/helpers"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

// RateLimiter func
func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(r, b)
	return func(c *gin.Context) {
		if !limiter.AllowN(time.Now(), 1) {
			c.JSON(http.StatusTooManyRequests, helpers.NewErrorResponse(http.StatusTooManyRequests, "Rate limit exceeded"))
			c.Abort()
			return
		}
		c.Next()
	}
}
