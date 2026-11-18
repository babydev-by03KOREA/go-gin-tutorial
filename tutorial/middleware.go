package tutorial

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func MyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		status := c.Writer.Status()
		fmt.Printf("상태코드: %d | 소요시간: %v\n", status, latency)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")

		if tokenString != "auth-token" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "인증되지 않은 사용자",
			})
			// 이후 기다리는 모든 Handler 코드를 중단하라!
			c.Abort()
			// 이 함수도 끝내라!
			return
		}
		c.Next()
	}
}
