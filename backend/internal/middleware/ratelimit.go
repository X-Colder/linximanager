package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimit 基于令牌桶的全局限流中间件
func RateLimit(qps float64, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(qps), burst)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    42900,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
