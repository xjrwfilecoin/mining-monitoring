package httpsvr

import (
	"fmt"
	"mining-monitoring/log"
	"mining-monitoring/model"
	"mining-monitoring/net/socket"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")
	//	c.Header("Access-Control-Expose-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}


func ListenAndServe(cfg *model.RuntimeConfig,server *socket.Server) {
	gin.SetMode(gin.ReleaseMode)
	httpRouter := gin.New()
	httpRouter.Use(cors())
	//httpRouter.Use(log.MyGinLogger(cfg.LogPath))
	httpRouter.Use(gin.Recovery())
	UseApiV1(httpRouter,server)
	httpSever := &http.Server{
		Addr:           cfg.HTTPListen,
		Handler:        httpRouter,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Info("web server start..." + time.Now().String())
	err := httpSever.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}
}
