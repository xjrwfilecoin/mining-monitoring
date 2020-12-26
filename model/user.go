package model

import "time"

type User struct {
	Id       string `bson:"_id" json:"id"`
	HeadPath string `json:"headPath" bson:"head_path"`
	PhoneNum string `json:"phoneNum" bson:"phone_num"`
	NickName string `json:"nickName" bson:"nick_name"` //昵称
	FullName string `json:"fullName" bson:"full_name"`
	Email    string `json:"email" bson:"email"`
}

type InnerUser struct {
	Id       string `json:"-" bson:"_id"`
	PhoneNum string `json:"-" bson:"phone_num"`
	Password string `json:"-" bson:"password"`
	Email    string `json:"-" bson:"email"`
}


type LoginLog struct {
	Id         string    `bson:"_id" json:"id"`
	Username   string    `json:"username" bson:"username"`
	Count      int64     `json:"count" bson:"count"`            // 指定时间内密码输入错误次数
	InsertTime int64     `json:"insertTime" bson:"insert_time"` // 记录插入的时间
	CreateTime time.Time `json:"createTime" bson:"create_time"`
}
