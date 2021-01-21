package statecode

var errorMap = map[int]string{
	Success:           "成功",
	UnknownError:      "未知错误",
	Fail:              "失败",
	TokenExpiredERROR: "token失效",
	ParamError:        "参数不对",


	PhoneNumNoExistError:  "手机号不存在",
	LoginMuchError:        "登录频率错误次数太多", //10001   登录次数太多
	AccountNotExistsError: "账号不存在",
	AccountOrPwdIsError:   "账号或者密码错误",
	EmailHaveBindErr:      "邮箱错误",
	PwdSameErr:            "新旧密码相同",
	PhoneHaveExistErr:     "手机号已经存在",
	GoogleAuthUnBindErr:   "未绑定google认证",
	HaveGoogleAuthErr:     "已经绑定google认证",
	OldPasswordErr:        "旧密码错误",
	IdCardNumHavExistErr:  "身份证号已经存在",
}

func CodeInfo(code int) string {
	str, ok := errorMap[code]
	if ok {
		return str
	}
	return errorMap[UnknownError]
}
