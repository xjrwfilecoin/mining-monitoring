package impl

import (
	"mining-monitoring/config"
	"mining-monitoring/model"
	"mining-monitoring/service"
	"mining-monitoring/statecode"
	"mining-monitoring/utils"
	"fmt"
	"time"
)

type UserService struct {

}

func (u *UserService) Login(username, password, devId string) (interface{}, int, error) {
	loginLog, err := u.loginLogDao.FindLoginLog(username)
	if err != nil {
		return nil, statecode.Fail, fmt.Errorf("get login log: %v \n", err)
	}
	ok := canLogin(loginLog)
	if !ok {
		return nil, statecode.LoginMuchError, fmt.Errorf("login much rate")
	}
	user, err := u.userDao.FindUser(username)
	if err != nil {
		return nil, statecode.Fail, fmt.Errorf("no find user %v \n", err)
	}
	if user.Id == "" {
		return nil, statecode.AccountNotExistsError, fmt.Errorf("user not exites ")
	}
	if user.Password != utils.MD5BySalt(password) {
		_, err := u.loginLogDao.UpdateLoginLog(username, loginLog.Count, loginLog.InsertTime)
		if err != nil {
			return nil, statecode.Fail, fmt.Errorf("update login log fail %v \n", err)
		}
		return nil, statecode.AccountOrPwdIsError, fmt.Errorf("username or password error  %v \n", username)
	}
	token, err := utils.GenerateTokenV1(user.Id, devId)
	if err != nil {
		return nil, statecode.Fail, err
	}
	_, err = u.userDao.UpdateToken(user.Id, token, devId)
	if err != nil {
		return nil, statecode.Fail, err
	}
	_, _ = u.loginLogDao.UpdateLoginLog(username, 0, 0)
	return token, statecode.Success, nil
}



func (u *UserService) PhoneLogin(username, code, devId string) (interface{}, int, error) {
	user, err := u.userDao.FindUser(username)
	if err != nil {
		return nil, statecode.Fail, err
	}
	if user.Id == "" {
		return nil, statecode.PhoneNumNoExistError, fmt.Errorf("user not exites ")
	}
	smsCode, err := u.smsDao.FindAndDeleteV1(username, code, config.PhoneLoginSmsCode)
	if err != nil {
		return nil, statecode.SmsCodeErr, err
	}
	if smsCode.IsExpired() {
		return nil, statecode.SmsCodeExpireErr, fmt.Errorf("sms code is expired \n")
	}
	token, err := utils.GenerateTokenV1(user.Id, devId)
	if err != nil {
		return nil, statecode.Fail, err
	}
	_, err = u.userDao.UpdateToken(user.Id, token, devId)
	if err != nil {
		return nil, statecode.Fail, err
	}
	return token, statecode.Success, nil
}




func NewUserService() service.IUserService {
	return &UserService{userDao: daoimpl.NewUserDao()}
}

func canLogin(loginLog *model.LoginLog) bool {
	if loginLog == nil || loginLog.Id == "" {
		return true
	}
	if time.Now().Unix()-loginLog.InsertTime < config.LoginErrorLimitIntervalTime &&
		loginLog.Count >= 5 {
		return false
	}
	return true
}