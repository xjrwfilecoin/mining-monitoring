package httpsvr

import (
	"github.com/gin-gonic/gin"
)

func UseApiV1(router *gin.Engine) {
	group := router.Group("/api/v1")
	group.POST("/socket.io/", nil)
}
