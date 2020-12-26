package service

// 用户相关
type IUserService interface {
	Login(username, password, devId string) (interface{}, int, error)
	PhoneLogin(username, code, devId string) (interface{}, int, error)
}

