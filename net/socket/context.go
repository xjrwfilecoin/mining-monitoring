package socket

import (
	"encoding/json"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"reflect"
)

type BaseFrom struct {
	Url   string      `json:"url"`
	Data  interface{} `json:"Body"`
	MsgId string      `json:"MsgId"`
}

func (b *BaseFrom) Valid() error {
	if b.Url != "" && b.MsgId != "" {
		return nil
	}
	return fmt.Errorf("form is error \n")
}

type ResponseCmd struct {
	Code    int         `json:"code"`
	Url     string      `json:"uri"`
	Message string      `json:"message"`
	Data    interface{} `json:"body"`
	MsgId   string      `json:"msgId"`
}

type IFromValid interface {
	Valid() error
}

type Context struct {
	Conn  socketio.Conn
	Body  string
	Uri   string
	MsgId string
}

func NewContext(conn socketio.Conn, uri, msgId, body string) *Context {
	return &Context{Conn: conn, Uri: uri, MsgId: msgId, Body: body}
}

func (c *Context) BindJson(obj IFromValid) error {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return fmt.Errorf("obj is not ptr \n")
	}
	if len(c.Body) == 0 {
		return fmt.Errorf("Body is empty \n")
	}
	err := json.Unmarshal([]byte(c.Body), obj)
	if err != nil {
		return err
	}
	err = obj.Valid()
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) FailResp(args ...string) {
	msg := "fail"
	if len(args) > 0 {
		msg = args[0]
	}
	resp := ResponseCmd{Code: 0, Url: c.Uri, Message: msg, Data: nil, MsgId: c.MsgId,}
	c.Conn.Emit(c.Uri, resp)
}

func (c *Context) SuccessResp(data interface{}) {
	resp := ResponseCmd{Code: 1, Url: c.Uri, Message: "success", Data: data, MsgId: c.MsgId,}
	c.Conn.Emit(c.Uri, resp)
}

func NewFailResp(msg ...string) *ResponseCmd {
	info := "fail"
	if len(msg) > 0 {
		info = msg[0]
	}
	return &ResponseCmd{Code: 0, Message: info}
}
