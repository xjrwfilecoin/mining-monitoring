package utils

import (
	"mining-monitoring/config"
	"crypto/tls"
	"github.com/go-gomail/gomail"
	"io/ioutil"
	"os"
	"strings"
)

// 发送邮件
func SendEmailLink(code string, email string) error {
	param := make(map[string]string)
	param["address"] = config.Email
	param["alias"] = "星际荣威官方"
	param["host"] = config.EmailHost
	param["password"] = config.EmailPassword
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(param["address"], param["alias"]))                                    //这种方式可以添加别名，即“XX官方”
	m.SetHeader("To", email)                                                                                  //发送给用户
	m.SetHeader("Subject", "星际荣威科技邮箱验证")                                                               //设置邮件主题
	m.SetBody("text/html", "星际荣威科技： <a href="+">请点击验证您的邮箱</a>,10分钟后失效") //设置邮件正文
	d := gomail.NewDialer(param["host"], 465, param["address"], param["password"])
	config := &tls.Config{InsecureSkipVerify: true}
	d.TLSConfig = config
	err := d.DialAndSend(m)
	return err

}

// 发送邮件
func SendEmailCode(code string, email string, operationType int) error {
	param := make(map[string]string)
	param["address"] = config.Email
	param["alias"] = "星际荣威官方"
	param["host"] = config.EmailHost
	param["password"] = config.EmailPassword
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(param["address"], param["alias"])) //这种方式可以添加别名，即“XX官方”
	m.SetHeader("To", email)                                               //发送给用户
	m.SetHeader("Subject", "星际荣威科技邮箱验证")                                   //设置邮件主题
	var templatePath = ""
	switch operationType {
	default:
		// todo
		templatePath = "./webroot/emailtemplate/index.html"

	}
	file, err := os.Open(templatePath)
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	body := strings.ReplaceAll(string(bytes), "XXXXX", code)
	m.SetBody("text/html", body) //设置邮件正文
	d := gomail.NewDialer(param["host"], 465, param["address"], param["password"])
	config := &tls.Config{InsecureSkipVerify: true}
	d.TLSConfig = config
	err = d.DialAndSend(m)
	return err
}
