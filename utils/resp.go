package utils

import (
	socketio "github.com/googollee/go-socket.io"
	"mining-monitoring/model"
)

func FailResp(url, msgId string, s socketio.Conn, args ...string) {
	msg := "fail"
	if len(args) > 0 {
		msg = args[0]
	}
	cmd := model.ResponseCmd{Code: 0, Url: url, Message: msg, Data: nil, MsgId: msgId,}
	s.Emit(url, cmd)
}

func SuccResp(url, msgId string, s socketio.Conn, data interface{}) {
	cmd := model.ResponseCmd{Code: 1, Url: url, Message: "success", Data: data, MsgId: msgId,}
	s.Emit(url, cmd)
}
