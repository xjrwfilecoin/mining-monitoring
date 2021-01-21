package utils

func IsEmpty(dest string) bool {
	if dest == "" || len(dest) < 1 {
		return true
	}
	return false
}

// 效验用户名，后期规范正则效验
func CheckUsernmae(username string) bool {
	if len(username) > 30 || len(username) < 6 {
		return false
	}
	return true
}

// 效验密码
func CheckPasswrod(password string) bool {
	if len(password) > 30 || len(password) < 6 {
		return false
	}
	return true
}
