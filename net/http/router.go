package httpsvr

import (
	"github.com/gin-gonic/gin"
	"mining-monitoring/net/socket"
)

func UseApiV1(router *gin.Engine,server *socket.Server) {
	router.GET("/socket.io/*any", gin.WrapH(server.GetServer()))
	router.POST("/socket.io/*any", gin.WrapH(server.GetServer()))
}
