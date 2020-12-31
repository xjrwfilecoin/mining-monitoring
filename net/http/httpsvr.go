package httpsvr

import (
	"fmt"
	"mining-monitoring/log"
	"mining-monitoring/model"
	"mining-monitoring/net/socket"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")

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
	httpRouter.Use(log.MyGinLogger(cfg.LogPath))
	httpRouter.Use(gin.Recovery())
	UseApiV1(httpRouter,server)
	// 静态资源目录
	webRootDir := "./webroot"
	if s, err := os.Stat(webRootDir); err != nil || !s.IsDir() {
		if err != nil {
			log.Logger.Fatalln("静态资源目录没创建...", err.Error())
		}
	}

	// 模板目录
	templatePath := "./webroot/templates/"
	if s, err := os.Stat(templatePath); err != nil || !s.IsDir() {
		if err != nil {
			log.Logger.Fatalln("html 模板目录没有创建...", err.Error())
		}
	}
	httpRouter.LoadHTMLGlob("./webroot/templates/*")
	httpRouter.StaticFS("/static", http.Dir(webRootDir))
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
