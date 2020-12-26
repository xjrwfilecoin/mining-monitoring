package model

import (
	"mining-monitoring/config"
	"time"
)

type SmsCode struct {
	Id         string    `bson:"_id" json:"id"`
	PhoneNum   string    `json:"phoneNum" bson:"phone_num"`
	Code       string    `json:"code"`
	OptionType int       `json:"optionType" bson:"option_type"` //验证码类型
	InsertTime int64     `json:"insertTime" bson:"insert_time"` //插入的时间的时间戳
	Status     string    `json:"status"`                        //验证码是否有效    1 有效，0 失效
	CreateTime time.Time `json:"createTime" bson:"create_time"`
}


func (code *SmsCode) IsExpired() bool{
	endTime := time.Now().Unix()
	if endTime-code.InsertTime > config.IntervalTime {
		return true
	}
	return false
}

// 邮箱验证码
type EmailCode struct {
	Id          string `bson:"_id" json:"id"`
	Uid         string             `json:"uid"`
	Email       string             `json:"email"`
	Code        string             `json:"code" bson:"code"`
	OptionType  int                `json:"optionType" bson:"option_type"`
	Status      int                `json:"status"`                          //邀请码状态
	ValidTime   int64              `json:"validTime" bson:"valid_time"`     //有效时间时长
	InvalidTime time.Time          `json:"invalidTime" bson:"invalid_time"` //失效时间
	CreateTime  time.Time          `json:"createTime" bson:"create_time"`
}

