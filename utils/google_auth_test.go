package utils

import (
	"fmt"
	"testing"
)

func TestGoogleAuth(t *testing.T) {

	fmt.Println("-----------------开启二次认证----------------------")
	user := "raojianli66@163.com"
	secret, code := initAuth(user)
	fmt.Println(secret, code)

	fmt.Println("-----------------信息校验----------------------")

	// secret最好持久化保存在
	// 验证,动态码(从谷歌验证器获取或者freeotp获取)
	bool, err := NewGoogleAuth().VerifyCode(secret, code)
	if bool {
		fmt.Println("√")
	} else {
		fmt.Println("X", err)
	}
}

func TestVerifyCode(t *testing.T) {
	b, e := NewGoogleAuth().VerifyCode("OTK2QSNL6FZUUU3MDKGSPUZVSMH7GI7I", "933868")
	fmt.Println(b, e)
}
