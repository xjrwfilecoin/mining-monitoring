package config


const FilecoinPrecision = uint64(1_000_000_000_000_000_000)

const SysTimefrom = "2006-01-02 15:04:05"

const LocationTimeZone = "Asia/Shanghai"

//邮箱地址
const EmailHost = "smtp.faidns.com"

// 二维码路径
const QRCodePath = "./webroot/qrcode/"

//  数据库没有找到信息错误
const NoResultErr = "mongo: no documents in result"

//限制用户登录的时间段
const LoginErrorLimitIntervalTime = 15 * 60 // 15分钟

// 验证码有效时间 单位 秒
const IntervalTime = 10 * 60 // 30分钟


// 上传文件限制大小 30M
const MaxUploadFileSize = 30 * 1024 * 1024