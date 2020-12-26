package utils

import "regexp"

var phoneReg = regexp.MustCompile(`^1[0-9]{10}$`)
var emailReg = regexp.MustCompile(`^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`)
var pwdReg = regexp.MustCompile(`^[a-z0-9A-Z_-]{6,18}`)
var usernameReg = regexp.MustCompile(`/^[a-z0-9_-]{3,16}$/`)
var ipReg = regexp.MustCompile(`/((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)/`)

var mongoIdReg = regexp.MustCompile(`^[a-z0-9A-Z]{24}$`)

// mongodb id校验
func VerifyMongoId(id string) bool {
	return mongoIdReg.MatchString(id)
}

// 手机号验证
func VerifyMobileFormat(mobileNum string) bool {
	return phoneReg.MatchString(mobileNum)
}

// 邮箱验证
func VerifyEmailFormat(email string) bool {
	return emailReg.MatchString(email)
}

//IP地址校验
func VerifyIp(ip string) bool {
	return ipReg.MatchString(ip)
}

// 用户名
func VerifyUserName(username string) bool {
	return usernameReg.MatchString(username)
}

//一个好用的检查密码强度的正则表达式,可以检查至少有一个大写,一个小写, 一个特殊字符,长度要是8:
func VerifyPassword(pwd string) bool {
	return pwdReg.MatchString(pwd)
}
