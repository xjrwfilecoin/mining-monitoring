package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"mining-monitoring/config"
	"time"
)

func GenerateTokenWithCode(uid string, code string, email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	// 设置token过期时间 todo 修改时间
	claims["exp"] = time.Now().Add(time.Second * 6000).Unix() //10分钟过期
	//时间戳
	claims["iat"] = time.Now().Unix()
	claims["uid"] = uid
	claims["code"] = code
	claims["email"] = email
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// verify token
func ValidTokenWithCode(tokenStr string) (string, string, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		//token验证逻辑,可以再次查数据验证用户是否存在
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		return "", "", "", err
	}
	if claim, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uid := claim["uid"]
		code := claim["code"]
		email := claim["email"]

		return uid.(string), code.(string), email.(string), nil
	} else {
		return "", "", "", errors.New("token verify fail")
	}

}

// get token
func GenerateToken(uid string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	// 设置token过期时间
	claims["exp"] = time.Now().Add(time.Hour * 5).Unix() // 2小时过期
	//时间戳
	claims["iat"] = time.Now().Unix()
	claims["uid"] = uid
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// verify token
func ValidToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		//token验证逻辑,可以再次查数据验证用户是否存在
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		return "", err
	}
	if claim, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uid := claim["uid"]
		return uid.(string), nil
	} else {
		return "", errors.New("token verify fail")
	}

}

// get token
func GenerateTokenV1(uid string, devId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	// 设置token过期时间
	claims["exp"] = time.Now().Add(time.Hour * 5).Unix() // 2小时过期
	//时间戳
	claims["iat"] = time.Now().Unix()
	claims["uid"] = uid
	claims["devId"] = devId
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// verify token
func ValidTokenV1(tokenStr string) (string, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		//token验证逻辑,可以再次查数据验证用户是否存在
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		return "", "", err
	}
	if claim, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uid := claim["uid"]
		devId := claim["devId"]
		return uid.(string), devId.(string), nil
	} else {
		return "", "", errors.New("token verify fail")
	}

}
