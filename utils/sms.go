package utils

import (
	"mining-monitoring/config"
	"mining-monitoring/log"
	"crypto/rand"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"io"
	"strings"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9'}

func GenVerifyCode() string {
	b := make([]byte, 6)
	n, err := io.ReadAtLeast(rand.Reader, b, 6)
	if n != 6 {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var smsClient *dysmsapi.Client
var err error

// send sms code
func SendSmsCode(phoneNum string, code string, optionType int) error {
	if smsClient == nil {
		smsClient, err = dysmsapi.NewClientWithAccessKey("cn-hangzhou", config.AliYunAccessKey, config.AliYunAccessKeySecret)
		if err != nil {
			return err
		}
	}
	// todo 短信模板从配置文件读取?
	request := dysmsapi.CreateSendSmsRequest()
	if optionType == 0 { // 注册
		request.TemplateCode = config.RegisterSmsTemplate
	}
	request.Scheme = "https"
	request.PhoneNumbers = phoneNum
	request.SignName = "星际荣威科技"
	request.TemplateParam = "{\"code\":" + code + "}"
	response, err := smsClient.SendSms(request)
	if err != nil {
		return err
	}
	if strings.EqualFold("OK", response.Code) {
		return nil
	} else {
		fmt.Println(response)
		log.Logger.Printf("sms send error: %s ", response)
		return fmt.Errorf("send sms fail %s", response)
	}
}
