package httpsvr

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mining-monitoring/utils"
	"net/http"
)

var DevIdError = errors.New("dev error")

func TokenVerifyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		uid, err := CheckToken(token)
		if err == DevIdError {
			c.JSON(http.StatusOK, gin.H{"code": -2, "msg": "认证失败"})
			c.Abort()
		} else if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "认证失败"})
			c.Abort()
		} else {
			c.Request.ParseForm()
			c.Request.PostForm.Add("uid", uid)
			c.Request.Form.Add("uid", uid)
			c.Next()
		}
	}
}

func CheckToken(token string) (string, error) {
	uid, err := utils.ValidToken(token)
	if err != nil {
		return uid, err
	}
	return uid, nil
}
