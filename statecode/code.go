package statecode

const (
	Success = iota + 1 // 1 成功
)

const (
	UnknownError      = iota - 1000 //-1000 未知错误
	Fail                            // -999 操作失败
	TokenExpiredERROR               // -998 token失效
	ParamError                      // -997 参数不对
)

const (
	PhoneNumNoExistError  = iota + 10000 //10000   该手机号未注册！
	LoginMuchError                       //10001   登录次数太多
	AccountNotExistsError                //10002   账号未注册，请注册账号！
	AccountOrPwdIsError                  //10003   账号密码错误
	EmailHaveBindErr                     //10004   邮箱已经绑定
	PwdSameErr                           //10005   密码想同
	PhoneHaveExistErr                    //10006   手机号存在
	GoogleAuthUnBindErr                  //10007    google认证未绑定
	HaveGoogleAuthErr                    //10008    存在google认证
	OldPasswordErr                       //10009    老密码错误
	IdCardNumHavExistErr                 //10011     身份证号已经存在
)

const (
	FileTooLargeError = iota + 20000 // 文件太大
	FileTooSmallError                // 文件太小
)

// 短信验证码
const (
	SmsCodeErr       = iota + 30000 // 30000 验证码错误
	SmsCodeExpireErr                // 30001 验证码失效
)

// 邮件验证码
const (
	EmailCodeErr       = iota + 40000 // email code错误
	EmailCodeExpireErr                // 验证码失效
)
