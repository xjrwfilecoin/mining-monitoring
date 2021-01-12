package model

type ResponseCmd struct {
	Code    int         `json:"code"`
	Url     string      `json:"url"`
	Message string      `json:"message"`
	Data    interface{} `json:"body"`
	MsgId   string      `json:"msgId"`
}