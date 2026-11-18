package tutorial

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func tutorialMain() {
	r := gin.Default()

	r.Use(MyLogger())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello GIN World",
		})
	})

	protected := r.Group("/admin")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"data": "관리자 권한 사용자만 보이는 메세지입니다.",
			})
		})
	}

	r.Run(":9090")
}
