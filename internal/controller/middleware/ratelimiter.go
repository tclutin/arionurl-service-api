package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
)

func RateLimiter() gin.HandlerFunc {
	limiter := rate.NewLimiter(1, 3)
	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "you have exceeded the number of requests per second"})
		}
	}
}
